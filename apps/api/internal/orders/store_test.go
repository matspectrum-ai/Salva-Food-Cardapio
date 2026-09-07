package orders

import "testing"

func TestMemoryStoreIdempotency(t *testing.T) {
	store := NewMemoryStore()
	order, _ := NewOrder("tenant", "order-1", "PDV", validItems())

	first, replay, err := store.Create(order, "key-1", "fp-1")
	if err != nil || replay {
		t.Fatalf("first create: replay=%v err=%v", replay, err)
	}
	second, replay, err := store.Create(order, "key-1", "fp-1")
	if err != nil || !replay {
		t.Fatalf("replay: replay=%v err=%v", replay, err)
	}
	if first.ID != second.ID {
		t.Fatalf("idempotent replay returned different order")
	}

	if _, _, err := store.Create(order, "key-1", "different"); err != ErrIdempotencyConflict {
		t.Fatalf("err = %v, want idempotency conflict", err)
	}
}

func TestMemoryStoreIsolatesTenants(t *testing.T) {
	store := NewMemoryStore()
	orderA, _ := NewOrder("tenant-a", "same-id", "PDV", validItems())
	orderB, _ := NewOrder("tenant-b", "same-id", "PDV", validItems())
	_, _, _ = store.Create(orderA, "key", "a")
	_, _, _ = store.Create(orderB, "key", "b")

	gotA, ok := store.Get("tenant-a", "same-id")
	if !ok || gotA.TenantID != "tenant-a" {
		t.Fatal("tenant A order missing")
	}
	gotB, ok := store.Get("tenant-b", "same-id")
	if !ok || gotB.TenantID != "tenant-b" {
		t.Fatal("tenant B order missing")
	}
}
