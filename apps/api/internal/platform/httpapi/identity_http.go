package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
)

func (a *API) registerIdentityRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", a.login)
	mux.HandleFunc("POST /api/v1/auth/logout", a.logout)
	mux.HandleFunc("GET /api/v1/me", a.me)
	mux.HandleFunc("GET /api/v1/collaborators", a.listCollaborators)
	mux.HandleFunc("POST /api/v1/collaborators", a.createCollaborator)
	mux.HandleFunc("PATCH /api/v1/collaborators/{id}/status", a.setCollaboratorStatus)
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(header) < 8 || !strings.EqualFold(header[:7], "Bearer ") {
		return ""
	}
	return strings.TrimSpace(header[7:])
}
func (a *API) principalFromRequest(w http.ResponseWriter, r *http.Request, requireActive bool) (identity.Principal, bool) {
	if a.authService == nil {
		writeError(w, http.StatusUnauthorized, errors.New("authentication is not configured"))
		return identity.Principal{}, false
	}
	token := bearerToken(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, errors.New("bearer token is required"))
		return identity.Principal{}, false
	}
	principal, err := a.authService.Authenticate(r.Context(), token)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, identity.ErrSessionNotFound) || errors.Is(err, identity.ErrSessionExpired) || errors.Is(err, identity.ErrSessionRevoked) {
			status = http.StatusUnauthorized
		}
		writeError(w, status, err)
		return identity.Principal{}, false
	}
	if requireActive && principal.Membership.Status != identity.MembershipActive {
		writeError(w, http.StatusForbidden, identity.ErrNoActiveMembership)
		return identity.Principal{}, false
	}
	return principal, true
}
func (a *API) tenantFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	if a.strictAuth {
		principal, ok := a.principalFromRequest(w, r, true)
		if !ok {
			return "", false
		}
		return principal.Membership.TenantID, true
	}
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, errors.New("X-Tenant-ID header is required"))
		return "", false
	}
	if _, err := uuid.Parse(tenantID); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("X-Tenant-ID must be a valid UUID"))
		return "", false
	}
	return tenantID, true
}

func requirePermission(w http.ResponseWriter, principal identity.Principal, permission identity.Permission) bool {
	if !principal.Membership.HasPermission(permission) {
		writeError(w, http.StatusForbidden, identity.ErrTenantAccessDenied)
		return false
	}
	return true
}
func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var input identity.LoginInput
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := a.authService.Login(r.Context(), input)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, identity.ErrInvalidCredentials):
			status = http.StatusUnauthorized
		case errors.Is(err, identity.ErrTenantSelectionRequired):
			status = http.StatusConflict
		case errors.Is(err, identity.ErrNoActiveMembership), errors.Is(err, identity.ErrTenantAccessDenied):
			status = http.StatusForbidden
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, errors.New("bearer token is required"))
		return
	}
	if err := a.authService.Logout(r.Context(), token); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, identity.ErrSessionNotFound) {
			status = http.StatusUnauthorized
		}
		writeError(w, status, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (a *API) me(w http.ResponseWriter, r *http.Request) {
	principal, ok := a.principalFromRequest(w, r, false)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, principal)
}

type collaboratorView struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Email       string                    `json:"email"`
	Phone       string                    `json:"phone"`
	ImageURL    string                    `json:"image_url,omitempty"`
	Title       string                    `json:"title"`
	Status      identity.MembershipStatus `json:"status"`
	Permissions []identity.Permission     `json:"permissions"`
}

func toCollaboratorView(value identity.Collaborator) collaboratorView {
	return collaboratorView{
		ID: value.User.ID, Name: value.User.Name, Email: value.User.Email,
		Phone: value.User.Phone, ImageURL: value.User.ImageURL,
		Title: value.Membership.Title, Status: value.Membership.Status,
		Permissions: append([]identity.Permission(nil), value.Membership.Permissions...),
	}
}
func (a *API) listCollaborators(w http.ResponseWriter, r *http.Request) {
	principal, ok := a.principalFromRequest(w, r, true)
	if !ok || !requirePermission(w, principal, identity.PermissionSettings) {
		return
	}
	items, err := a.identityService.ListCollaborators(r.Context(), principal.Membership.TenantID, r.URL.Query().Get("search"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	result := make([]collaboratorView, 0, len(items))
	for _, item := range items {
		result = append(result, toCollaboratorView(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": result})
}

func (a *API) createCollaborator(w http.ResponseWriter, r *http.Request) {
	principal, ok := a.principalFromRequest(w, r, true)
	if !ok || !requirePermission(w, principal, identity.PermissionSettings) {
		return
	}
	var input identity.CreateCollaboratorInput
	if !decodeJSON(w, r, &input) {
		return
	}
	created, err := a.identityService.CreateCollaborator(r.Context(), principal.Membership.TenantID, input)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, identity.ErrUserAlreadyExists) || errors.Is(err, identity.ErrMembershipAlreadyExists) {
			status = http.StatusConflict
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusCreated, toCollaboratorView(created))
}

type setCollaboratorStatusRequest struct {
	Status identity.MembershipStatus `json:"status"`
}

func (a *API) setCollaboratorStatus(w http.ResponseWriter, r *http.Request) {
	principal, ok := a.principalFromRequest(w, r, true)
	if !ok || !requirePermission(w, principal, identity.PermissionSettings) {
		return
	}
	var input setCollaboratorStatusRequest
	if !decodeJSON(w, r, &input) {
		return
	}
	if _, err := uuid.Parse(r.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("collaborator id must be a valid UUID"))
		return
	}
	if err := a.identityService.SetMembershipStatus(r.Context(), principal.Membership.TenantID, r.PathValue("id"), input.Status); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, identity.ErrCollaboratorNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
