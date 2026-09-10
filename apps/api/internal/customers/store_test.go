package customers

import (
	"context"
	"testing"
)

func TestMemoryStoreCustomerAndAddressIsolation(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()
	c, err := NewCustomer("t1", "c1", "Ana", "91999999999", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateCustomer(ctx, c); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateCustomer(ctx, c); err != ErrAlreadyExists {
		t.Fatal(err)
	}

	a, err := NewAddress("t1", "a1", "c1", "Rua A", "10", "Centro", "Santarém", "PA", "68000-000")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateAddress(ctx, a); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := s.GetCustomer(ctx, "t2", "c1"); ok {
		t.Fatal("cross-tenant customer visible")
	}
	if _, ok, _ := s.GetAddress(ctx, "t2", "a1"); ok {
		t.Fatal("cross-tenant address visible")
	}
}

func TestMemoryStoreAddressRequiresExistingCustomer(t *testing.T) {
	s := NewMemoryStore()
	a, _ := NewAddress("t1", "a1", "missing", "Rua A", "10", "Centro", "Santarém", "PA", "68000-000")
	if err := s.CreateAddress(context.Background(), a); err != ErrCustomerNotFound {
		t.Fatal(err)
	}
}
