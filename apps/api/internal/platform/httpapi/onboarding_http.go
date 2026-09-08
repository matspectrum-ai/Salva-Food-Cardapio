package httpapi

import (
	"errors"
	"net/http"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/tenant"
)

func (a *API) registerOnboardingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/onboarding", a.createOnboarding)
}

type onboardingResponse struct {
	Tenant        tenant.Tenant        `json:"tenant"`
	Establishment tenant.Establishment `json:"establishment"`
	Owner         collaboratorView     `json:"owner"`
}

func (a *API) createOnboarding(w http.ResponseWriter, r *http.Request) {
	var input onboarding.Input
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := a.onboardingService.Create(r.Context(), input)
	if err != nil {
		writeError(w, onboardingStatus(err), err)
		return
	}
	writeJSON(w, http.StatusCreated, onboardingResponse{
		Tenant:        result.Tenant,
		Establishment: result.Establishment,
		Owner:         toCollaboratorView(result.Owner),
	})
}

func onboardingStatus(err error) int {
	switch {
	case errors.Is(err, onboarding.ErrTenantAlreadyExists),
		errors.Is(err, onboarding.ErrEstablishmentAlreadyExists),
		errors.Is(err, identity.ErrUserAlreadyExists),
		errors.Is(err, identity.ErrMembershipAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, tenant.ErrNameRequired),
		errors.Is(err, tenant.ErrInvalidTimezone),
		errors.Is(err, identity.ErrNameRequired),
		errors.Is(err, identity.ErrInvalidCPF),
		errors.Is(err, identity.ErrInvalidEmail),
		errors.Is(err, identity.ErrPhoneRequired),
		errors.Is(err, identity.ErrWeakPassword):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
