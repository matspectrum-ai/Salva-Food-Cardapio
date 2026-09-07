package orders

import (
	"errors"
	"testing"
)

type fakeCatalog struct {
	products map[string]CatalogProduct
	err      error
	calls    int
}

func (f *fakeCatalog) LookupOrderProduct(_, itemID string) (CatalogProduct, error) {
	f.calls++
	if f.err != nil {
		return CatalogProduct{}, f.err
	}
	product, ok := f.products[itemID]
	if !ok {
		return CatalogProduct{}, errors.New("not found")
	}
	return product, nil
}

func TestServiceCreatesSnapshotFromCatalog(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{
		"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000},
	}}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" })
	order, replay, err := service.Create("tenant", "idem-1", CreateInput{
		Source: "PDV",
		Items:  []CreateLineInput{{ItemID: "item-1", Quantity: 2}},
	})
	if err != nil || replay {
		t.Fatalf("create: replay=%v err=%v", replay, err)
	}
	if order.TotalCents != 6000 || order.Items[0].Name != "Pizza" {
		t.Fatalf("unexpected snapshot: %+v", order)
	}
}
func TestServiceReplaysBeforeReadingCatalogAgain(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{
		"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000},
	}}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" })
	input := CreateInput{Source: "PDV", Items: []CreateLineInput{{ItemID: "item-1", Quantity: 1}}}

	first, _, err := service.Create("tenant", "idem-1", input)
	if err != nil {
		t.Fatal(err)
	}
	catalog.products["item-1"] = CatalogProduct{ID: "item-1", Name: "Pizza", PriceCents: 9999, SoldOut: true}
	second, replay, err := service.Create("tenant", "idem-1", input)
	if err != nil || !replay {
		t.Fatalf("replay=%v err=%v", replay, err)
	}
	if catalog.calls != 1 {
		t.Fatalf("catalog calls = %d, want 1", catalog.calls)
	}
	if second.ID != first.ID || second.TotalCents != 3000 {
		t.Fatalf("replay changed snapshot: %+v", second)
	}
}

func TestServiceRejectsUnavailableProductAndMissingIdempotencyKey(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{
		"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000, SoldOut: true},
	}}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" })
	input := CreateInput{Source: "PDV", Items: []CreateLineInput{{ItemID: "item-1", Quantity: 1}}}

	if _, _, err := service.Create("tenant", "", input); err != ErrIdempotencyKeyRequired {
		t.Fatalf("err = %v, want missing idempotency key", err)
	}
	if _, _, err := service.Create("tenant", "idem-1", input); err != ErrCatalogItemUnavailable {
		t.Fatalf("err = %v, want unavailable", err)
	}
}
