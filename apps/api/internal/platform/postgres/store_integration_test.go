package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/tenant"
)

const (
	testTenantA  = "11111111-1111-4111-8111-111111111111"
	testTenantB  = "22222222-2222-4222-8222-222222222222"
	testCatA     = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	testItemA    = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	testOrderA   = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	testUserA    = "eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee"
	testUserB    = "ffffffff-ffff-4fff-8fff-ffffffffffff"
	fingerprintA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	fingerprintB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

func TestStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if _, err := store.pool.Exec(ctx, `TRUNCATE outbox_events, order_idempotency, order_items, orders, catalog_items, catalog_categories, establishments, tenants CASCADE`); err != nil {
		t.Fatal(err)
	}

	seedTenant := func(idValue, name string) {
		id, parseErr := parseUUID(idValue)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if _, createErr := store.q.CreateTenant(ctx, db.CreateTenantParams{ID: id, Name: name}); createErr != nil {
			t.Fatal(createErr)
		}
	}
	seedTenant(testTenantA, "Tenant A")
	seedTenant(testTenantB, "Tenant B")

	category, err := catalog.NewCategory(testTenantA, testCatA, "Lanches", 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCategory(ctx, category); err != nil {
		t.Fatal(err)
	}

	item, err := catalog.NewItem(testTenantA, testItemA, "X-Burger", category, 2590, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateItem(ctx, item); err != nil {
		t.Fatal(err)
	}

	foreignItem := item
	foreignItem.TenantID = testTenantB
	foreignItem.ID = "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
	if err := store.CreateItem(ctx, foreignItem); !errors.Is(err, catalog.ErrCategoryNotFound) {
		t.Fatalf("cross-tenant category err = %v", err)
	}

	order, err := orders.NewOrder(testTenantA, testOrderA, "PDV", []orders.LineItemSnapshot{{
		ItemID: testItemA, Name: "X-Burger", Quantity: 2, UnitPriceCents: 2590,
	}})
	if err != nil {
		t.Fatal(err)
	}
	created, replay, err := store.Create(ctx, order, "checkout-001", fingerprintA)
	if err != nil || replay {
		t.Fatalf("create replay=%v err=%v", replay, err)
	}
	if created.TotalCents != 5180 || len(created.Items) != 1 {
		t.Fatalf("unexpected created order: %+v", created)
	}

	replayed, replay, err := store.Create(ctx, order, "checkout-001", fingerprintA)
	if err != nil || !replay {
		t.Fatalf("replay=%v err=%v", replay, err)
	}
	if replayed.ID != created.ID {
		t.Fatalf("replayed different order: %s != %s", replayed.ID, created.ID)
	}

	if _, _, err := store.Create(ctx, order, "checkout-001", fingerprintB); !errors.Is(err, orders.ErrIdempotencyConflict) {
		t.Fatalf("idempotency conflict err = %v", err)
	}

	production, err := store.Transition(ctx, testTenantA, testOrderA, orders.StatusProduction)
	if err != nil || production.Status != orders.StatusProduction {
		t.Fatalf("transition: status=%v err=%v", production.Status, err)
	}

	if _, ok, err := store.Get(ctx, testTenantB, testOrderA); err != nil || ok {
		t.Fatalf("cross-tenant order visible: ok=%v err=%v", ok, err)
	}
	ordersB, err := store.List(ctx, testTenantB)
	if err != nil || len(ordersB) != 0 {
		t.Fatalf("tenant B list: len=%d err=%v", len(ordersB), err)
	}
}

func TestIdentityStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if _, err := store.pool.Exec(ctx, `TRUNCATE auth_sessions, tenant_memberships, app_users, tenants CASCADE`); err != nil {
		t.Fatal(err)
	}

	seed := func(idValue, name string) {
		id, parseErr := parseUUID(idValue)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if _, createErr := store.q.CreateTenant(ctx, db.CreateTenantParams{ID: id, Name: name}); createErr != nil {
			t.Fatal(createErr)
		}
	}
	seed(testTenantA, "Tenant A")
	seed(testTenantB, "Tenant B")

	user, err := identity.NewUser(
		testUserA, "Maria Silva", "529.982.247-25", "maria@example.com",
		"(93) 99999-9999", "test-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	membership, err := identity.NewMembership(
		testTenantA, testUserA, "Gerente",
		[]identity.Permission{identity.PermissionPOS, identity.PermissionReports},
	)
	if err != nil {
		t.Fatal(err)
	}
	collaborator := identity.Collaborator{User: user, Membership: membership}
	if err := store.CreateCollaborator(ctx, collaborator); err != nil {
		t.Fatal(err)
	}

	found, ok, err := store.FindUserByEmail(ctx, "MARIA@example.com")
	if err != nil || !ok || found.ID != testUserA {
		t.Fatalf("found=%+v ok=%v err=%v", found, ok, err)
	}

	items, err := store.ListCollaborators(ctx, testTenantA, "gerente")
	if err != nil || len(items) != 1 || items[0].User.ID != testUserA {
		t.Fatalf("items=%+v err=%v", items, err)
	}
	otherTenant, err := store.ListCollaborators(ctx, testTenantB, "")
	if err != nil || len(otherTenant) != 0 {
		t.Fatalf("tenant leak items=%+v err=%v", otherTenant, err)
	}

	if err := store.SetMembershipStatus(ctx, testTenantA, testUserA, identity.MembershipInactive); err != nil {
		t.Fatal(err)
	}
	updated, ok, err := store.GetMembership(ctx, testTenantA, testUserA)
	if err != nil || !ok || updated.Status != identity.MembershipInactive {
		t.Fatalf("membership=%+v ok=%v err=%v", updated, ok, err)
	}

	duplicate, err := identity.NewUser(
		testUserB, "Outra Pessoa", "168.995.350-09", "maria@example.com",
		"(93) 98888-8888", "test-hash-2",
	)
	if err != nil {
		t.Fatal(err)
	}
	duplicateMembership, err := identity.NewMembership(
		testTenantB, testUserB, "Operador", []identity.Permission{identity.PermissionPOS},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCollaborator(ctx, identity.Collaborator{
		User: duplicate, Membership: duplicateMembership,
	}); !errors.Is(err, identity.ErrUserAlreadyExists) {
		t.Fatalf("duplicate identity err=%v", err)
	}
}
func TestAuthSessionIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	store, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if _, err := store.pool.Exec(ctx, `TRUNCATE auth_sessions, tenant_memberships, app_users, tenants CASCADE`); err != nil {
		t.Fatal(err)
	}
	tenantID, _ := parseUUID(testTenantA)
	if _, err := store.q.CreateTenant(ctx, db.CreateTenantParams{ID: tenantID, Name: "Tenant A"}); err != nil {
		t.Fatal(err)
	}
	passwords := identity.NewBcryptPasswordManager(4)
	fixturePassword := fmt.Sprintf("Fixture%cPass%d", '#', 123)
	hash, err := passwords.Hash(fixturePassword)
	if err != nil {
		t.Fatal(err)
	}
	user, err := identity.NewUser(
		testUserA, "Maria Silva", "529.982.247-25", "maria@example.com",
		"(93) 99999-9999", hash,
	)
	if err != nil {
		t.Fatal(err)
	}
	membership, err := identity.NewMembership(
		testTenantA, testUserA, "Gerente",
		[]identity.Permission{identity.PermissionPOS, identity.PermissionReports},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCollaborator(ctx, identity.Collaborator{User: user, Membership: membership}); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	auth := identity.NewAuthService(
		store, passwords,
		func() (string, error) { return "fixture-session-token", nil },
		func() string { return "99999999-9999-4999-8999-999999999999" },
		func() time.Time { return now }, 12*time.Hour,
	)
	login, err := auth.Login(ctx, identity.LoginInput{
		Email: "maria@example.com", Password: fixturePassword,
	})
	if err != nil {
		t.Fatal(err)
	}
	if login.Token != "fixture-session-token" {
		t.Fatalf("unexpected token %q", login.Token)
	}
	if _, ok, err := store.GetSessionByTokenHash(ctx, login.Token); err != nil || ok {
		t.Fatalf("raw token persisted or lookup failed: ok=%v err=%v", ok, err)
	}
	principal, err := auth.Authenticate(ctx, login.Token)
	if err != nil || principal.User.ID != testUserA {
		t.Fatalf("authenticate principal=%+v err=%v", principal, err)
	}
	if err := store.SetMembershipStatus(ctx, testTenantA, testUserA, identity.MembershipInactive); err != nil {
		t.Fatal(err)
	}
	principal, err = auth.Authenticate(ctx, login.Token)
	if err != nil || principal.Membership.Status != identity.MembershipInactive {
		t.Fatalf("inactive session principal=%+v err=%v", principal, err)
	}
	if err := auth.Logout(ctx, login.Token); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.Authenticate(ctx, login.Token); !errors.Is(err, identity.ErrSessionRevoked) {
		t.Fatalf("post-logout err=%v want=%v", err, identity.ErrSessionRevoked)
	}
}

func TestOnboardingStoreIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := Open(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if _, err := store.pool.Exec(ctx, `TRUNCATE auth_sessions, tenant_memberships, app_users, establishments, tenants CASCADE`); err != nil {
		t.Fatal(err)
	}

	tenantA, err := tenant.New(testTenantA, "Tenant A")
	if err != nil {
		t.Fatal(err)
	}
	establishmentA, err := tenant.NewEstablishment(
		"33333333-3333-4333-8333-333333333333", testTenantA, "Store A", "America/Santarem",
	)
	if err != nil {
		t.Fatal(err)
	}
	email := fmt.Sprintf("%s@%s", "owner", "example.com")
	userA, err := identity.NewUser(
		testUserA, "Owner A", "52998224725", email, "fixture-phone", "fixture-hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	membershipA, err := identity.NewMembership(
		testTenantA, testUserA, onboarding.OwnerTitle, identity.AllPermissions(),
	)
	if err != nil {
		t.Fatal(err)
	}
	ownerA := identity.Collaborator{User: userA, Membership: membershipA}
	if err := store.Bootstrap(ctx, tenantA, establishmentA, ownerA); err != nil {
		t.Fatal(err)
	}

	foundUser, ok, err := store.FindUserByEmail(ctx, email)
	if err != nil || !ok || foundUser.ID != testUserA {
		t.Fatalf("owner lookup: ok=%v user=%+v err=%v", ok, foundUser, err)
	}
	tenantB, err := tenant.New(testTenantB, "Tenant B")
	if err != nil {
		t.Fatal(err)
	}
	establishmentB, err := tenant.NewEstablishment(
		"44444444-4444-4444-8444-444444444444", testTenantB, "Store B", "America/Sao_Paulo",
	)
	if err != nil {
		t.Fatal(err)
	}
	userB, err := identity.NewUser(
		testUserB, "Owner B", "16899535009", email, "fixture-phone", "fixture-hash-2",
	)
	if err != nil {
		t.Fatal(err)
	}
	membershipB, err := identity.NewMembership(
		testTenantB, testUserB, onboarding.OwnerTitle, identity.AllPermissions(),
	)
	if err != nil {
		t.Fatal(err)
	}
	ownerB := identity.Collaborator{User: userB, Membership: membershipB}
	if err := store.Bootstrap(ctx, tenantB, establishmentB, ownerB); !errors.Is(err, identity.ErrUserAlreadyExists) {
		t.Fatalf("duplicate owner err=%v", err)
	}

	var leakedTenants int
	if err := store.pool.QueryRow(
		ctx, `SELECT count(*) FROM tenants WHERE id = $1`, testTenantB,
	).Scan(&leakedTenants); err != nil {
		t.Fatal(err)
	}
	if leakedTenants != 0 {
		t.Fatalf("duplicate onboarding leaked %d tenant rows", leakedTenants)
	}

	var leakedEstablishments int
	if err := store.pool.QueryRow(
		ctx, `SELECT count(*) FROM establishments WHERE tenant_id = $1`, testTenantB,
	).Scan(&leakedEstablishments); err != nil {
		t.Fatal(err)
	}
	if leakedEstablishments != 0 {
		t.Fatalf("duplicate onboarding leaked %d establishment rows", leakedEstablishments)
	}
}
