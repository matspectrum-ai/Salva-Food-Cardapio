package orders

import (
	"context"
	"testing"
)

func TestMemoryStoreIdempotency(t *testing.T) {
	store := NewMemoryStore()
	order, _ := NewOrder("tenant", "order-1", "PDV", validItems())

	first, replay, err := store.Create(context.Background(), order, "key-1", "fp-1")
	if err != nil || replay {
		t.Fatalf("first create: replay=%v err=%v", replay, err)
	}
	second, replay, err := store.Create(context.Background(), order, "key-1", "fp-1")
	if err != nil || !replay {
		t.Fatalf("replay: replay=%v err=%v", replay, err)
	}
	if first.ID != second.ID {
		t.Fatalf("idempotent replay returned different order")
	}

	if _, _, err := store.Create(context.Background(), order, "key-1", "different"); err != ErrIdempotencyConflict {
		t.Fatalf("err = %v, want idempotency conflict", err)
	}
}

func TestMemoryStoreIsolatesTenants(t *testing.T) {
	store := NewMemoryStore()
	orderA, _ := NewOrder("tenant-a", "same-id", "PDV", validItems())
	orderB, _ := NewOrder("tenant-b", "same-id", "PDV", validItems())
	_, _, _ = store.Create(context.Background(), orderA, "key", "a")
	_, _, _ = store.Create(context.Background(), orderB, "key", "b")

	gotA, ok, err := store.Get(context.Background(), "tenant-a", "same-id")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || gotA.TenantID != "tenant-a" {
		t.Fatal("tenant A order missing")
	}
	gotB, ok, err := store.Get(context.Background(), "tenant-b", "same-id")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || gotB.TenantID != "tenant-b" {
		t.Fatal("tenant B order missing")
	}
}
