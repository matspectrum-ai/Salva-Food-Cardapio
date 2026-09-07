package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
)

const (
	testTenantA  = "11111111-1111-4111-8111-111111111111"
	testTenantB  = "22222222-2222-4222-8222-222222222222"
	testCatA     = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	testItemA    = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	testOrderA   = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
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
