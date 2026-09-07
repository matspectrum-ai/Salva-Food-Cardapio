package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
)

func authFixturePassword() string {
	return fmt.Sprintf("Fixture%cPass%d", '#', 123)
}

type authenticatedHarness struct {
	testHarness
	identityStore *identity.MemoryStore
}

func newAuthenticatedHarness(t *testing.T) authenticatedHarness {
	t.Helper()
	catalogStore := catalog.NewMemoryStore()
	orderStore := orders.NewMemoryStore()
	identityStore := identity.NewMemoryStore()
	var sequence atomic.Int64
	newID := func() string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", sequence.Add(1)) }
	passwords := identity.NewBcryptPasswordManager(4)
	identityService := identity.NewService(identityStore, passwords, newID)
	if _, err := identityService.CreateCollaborator(context.Background(), tenantA, identity.CreateCollaboratorInput{
		Name: "Owner", CPF: "529.982.247-25", Email: "owner@example.com",
		Password: authFixturePassword(), Phone: "(93) 99999-9999", Title: "Owner",
		Permissions: []identity.Permission{identity.PermissionSettings, identity.PermissionPOS, identity.PermissionReports},
	}); err != nil {
		t.Fatal(err)
	}
	var tokenSequence atomic.Int64
	authService := identity.NewAuthService(
		identityStore, passwords,
		func() (string, error) { return fmt.Sprintf("fixture-token-%d", tokenSequence.Add(1)), nil },
		newID, func() time.Time { return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC) }, 12*time.Hour,
	)
	api := NewAuthenticated(catalogStore, orderStore, identityService, authService, newID)
	return authenticatedHarness{testHarness: testHarness{handler: api.Handler()}, identityStore: identityStore}
}
func loginOwner(t *testing.T, h authenticatedHarness) string {
	t.Helper()
	rec := h.do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email": "owner@example.com", "password": authFixturePassword(),
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: got %d body=%s", rec.Code, rec.Body.String())
	}
	result := decodeResponse[identity.LoginResult](t, rec)
	if result.Token == "" || result.Principal.Membership.TenantID != tenantA {
		t.Fatalf("unexpected login result: %+v", result)
	}
	return result.Token
}

func TestAuthenticatedTenantComesFromSession(t *testing.T) {
	h := newAuthenticatedHarness(t)
	token := loginOwner(t, h)
	headers := map[string]string{"Authorization": "Bearer " + token}
	created := h.do(t, http.MethodPost, "/api/v1/catalog/categories", tenantB, map[string]any{
		"name": "Bebidas", "sort_order": 0,
	}, headers)
	if created.Code != http.StatusCreated {
		t.Fatalf("create category: got %d body=%s", created.Code, created.Body.String())
	}
	category := decodeResponse[catalog.Category](t, created)
	if category.TenantID != tenantA {
		t.Fatalf("tenant came from spoofed header: %s", category.TenantID)
	}
}
func TestCollaboratorHTTPAndPermissionGate(t *testing.T) {
	h := newAuthenticatedHarness(t)
	ownerToken := loginOwner(t, h)
	headers := map[string]string{"Authorization": "Bearer " + ownerToken}

	created := h.do(t, http.MethodPost, "/api/v1/collaborators", "", map[string]any{
		"name": "Operador", "cpf": "168.995.350-09", "email": "operator@example.com",
		"password": authFixturePassword(), "phone": "(93) 98888-8888", "title": "Operador",
		"permissions": []string{string(identity.PermissionPOS)},
	}, headers)
	if created.Code != http.StatusCreated {
		t.Fatalf("create collaborator: got %d body=%s", created.Code, created.Body.String())
	}
	operator := decodeResponse[collaboratorView](t, created)
	if operator.ID == "" || operator.Title != "Operador" {
		t.Fatalf("unexpected collaborator: %+v", operator)
	}

	listed := h.do(t, http.MethodGet, "/api/v1/collaborators?search=operador", "", nil, headers)
	if listed.Code != http.StatusOK {
		t.Fatalf("list collaborators: got %d body=%s", listed.Code, listed.Body.String())
	}
	operatorLogin := h.do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email": "operator@example.com", "password": authFixturePassword(),
	}, nil)
	if operatorLogin.Code != http.StatusOK {
		t.Fatalf("operator login: got %d body=%s", operatorLogin.Code, operatorLogin.Body.String())
	}
	operatorToken := decodeResponse[identity.LoginResult](t, operatorLogin).Token
	forbidden := h.do(t, http.MethodGet, "/api/v1/collaborators", "", nil, map[string]string{
		"Authorization": "Bearer " + operatorToken,
	})
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("permission gate: got %d body=%s", forbidden.Code, forbidden.Body.String())
	}

	statusRec := h.do(t, http.MethodPatch, "/api/v1/collaborators/"+operator.ID+"/status", "", map[string]any{
		"status": identity.MembershipInactive,
	}, headers)
	if statusRec.Code != http.StatusNoContent {
		t.Fatalf("disable collaborator: got %d body=%s", statusRec.Code, statusRec.Body.String())
	}
	blockedLogin := h.do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email": "operator@example.com", "password": authFixturePassword(),
	}, nil)
	if blockedLogin.Code != http.StatusForbidden {
		t.Fatalf("inactive login: got %d body=%s", blockedLogin.Code, blockedLogin.Body.String())
	}
}
func TestMeLogoutAndStrictAuth(t *testing.T) {
	h := newAuthenticatedHarness(t)
	missing := h.do(t, http.MethodGet, "/api/v1/orders", tenantA, nil, nil)
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("strict auth: got %d body=%s", missing.Code, missing.Body.String())
	}
	token := loginOwner(t, h)
	headers := map[string]string{"Authorization": "Bearer " + token}
	me := h.do(t, http.MethodGet, "/api/v1/me", "", nil, headers)
	if me.Code != http.StatusOK {
		t.Fatalf("me: got %d body=%s", me.Code, me.Body.String())
	}
	logout := h.do(t, http.MethodPost, "/api/v1/auth/logout", "", nil, headers)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout: got %d body=%s", logout.Code, logout.Body.String())
	}
	revoked := h.do(t, http.MethodGet, "/api/v1/me", "", nil, headers)
	if revoked.Code != http.StatusUnauthorized {
		t.Fatalf("revoked session: got %d body=%s", revoked.Code, revoked.Body.String())
	}
}
