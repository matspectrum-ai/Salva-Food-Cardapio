package httpapi

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
)

func onboardingFixturePassword() string {
	return fmt.Sprintf("Fixture%cOnboard%d", '#', 123)
}

func newOnboardingHarness() testHarness {
	catalogStore := catalog.NewMemoryStore()
	orderStore := orders.NewMemoryStore()
	identityStore := identity.NewMemoryStore()
	onboardingStore := onboarding.NewMemoryRepository(identityStore)
	passwords := identity.NewBcryptPasswordManager(4)

	var sequence atomic.Int64
	newID := func() string {
		return fmt.Sprintf("00000000-0000-4000-8000-%012d", sequence.Add(1))
	}
	identityService := identity.NewService(identityStore, passwords, newID)
	authService := identity.NewAuthService(
		identityStore,
		passwords,
		func() (string, error) { return "fixture-onboarding-token", nil },
		newID,
		func() time.Time { return time.Date(2026, 9, 7, 18, 0, 0, 0, time.UTC) },
		12*time.Hour,
	)
	onboardingService := onboarding.NewService(onboardingStore, passwords, newID)
	api := NewAuthenticatedWithOnboarding(
		catalogStore, orderStore, identityService, authService, onboardingService, newID,
	)
	return testHarness{handler: api.Handler()}
}

func TestOnboardingHTTPCreatesOwnerAndAllowsLogin(t *testing.T) {
	h := newOnboardingHarness()
	email := fmt.Sprintf("%s@%s", "owner", "example.com")
	body := map[string]any{
		"tenant_name":        "Dinda Foods",
		"establishment_name": "Dinda Foods Centro",
		"timezone":           "America/Santarem",
		"owner": map[string]any{
			"name":     "Owner Test",
			"cpf":      "52998224725",
			"email":    email,
			"phone":    "fixture-phone",
			"password": onboardingFixturePassword(),
		},
	}
	rec := h.do(t, http.MethodPost, "/api/v1/onboarding", "", body, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("onboarding: got %d body=%s", rec.Code, rec.Body.String())
	}
	result := decodeResponse[onboardingResponse](t, rec)
	if result.Tenant.ID == "" || result.Establishment.TenantID != result.Tenant.ID {
		t.Fatalf("unexpected onboarding result: %+v", result)
	}
	if result.Owner.Title != onboarding.OwnerTitle || len(result.Owner.Permissions) != len(identity.AllPermissions()) {
		t.Fatalf("owner view=%+v", result.Owner)
	}
	if strings.Contains(rec.Body.String(), `"cpf"`) || strings.Contains(rec.Body.String(), "password_hash") {
		t.Fatalf("sensitive identity data leaked: %s", rec.Body.String())
	}

	login := h.do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": onboardingFixturePassword(),
	}, nil)
	if login.Code != http.StatusOK {
		t.Fatalf("owner login: got %d body=%s", login.Code, login.Body.String())
	}
	loginResult := decodeResponse[identity.LoginResult](t, login)
	if loginResult.Principal.Membership.TenantID != result.Tenant.ID {
		t.Fatalf("login tenant=%s want=%s", loginResult.Principal.Membership.TenantID, result.Tenant.ID)
	}

	body["tenant_name"] = "Outro Tenant"
	body["establishment_name"] = "Outra Unidade"
	body["owner"].(map[string]any)["cpf"] = "16899535009"
	duplicate := h.do(t, http.MethodPost, "/api/v1/onboarding", "", body, nil)
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate onboarding: got %d body=%s", duplicate.Code, duplicate.Body.String())
	}
}

func TestOnboardingHTTPRejectsInvalidTimezone(t *testing.T) {
	h := newOnboardingHarness()
	rec := h.do(t, http.MethodPost, "/api/v1/onboarding", "", map[string]any{
		"tenant_name":        "Tenant",
		"establishment_name": "Store",
		"timezone":           "Invalid/Timezone",
		"owner": map[string]any{
			"name":     "Owner",
			"cpf":      "52998224725",
			"email":    fmt.Sprintf("%s@%s", "timezone", "example.com"),
			"phone":    "fixture-phone",
			"password": onboardingFixturePassword(),
		},
	}, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid timezone: got %d body=%s", rec.Code, rec.Body.String())
	}
}
