package catalog

import "testing"

func TestMemoryStoreIsolatesTenants(t *testing.T) {
	store := NewMemoryStore()
	catA, _ := NewCategory("tenant-a", "same-id", "A", 0)
	catB, _ := NewCategory("tenant-b", "same-id", "B", 0)
	if err := store.CreateCategory(catA); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCategory(catB); err != nil {
		t.Fatal(err)
	}

	gotA, ok := store.GetCategory("tenant-a", "same-id")
	if !ok || gotA.Name != "A" {
		t.Fatalf("tenant A leaked or missing: %+v %v", gotA, ok)
	}
	gotB, ok := store.GetCategory("tenant-b", "same-id")
	if !ok || gotB.Name != "B" {
		t.Fatalf("tenant B leaked or missing: %+v %v", gotB, ok)
	}
}

func TestMemoryStoreListsCatalogInExplicitOrder(t *testing.T) {
	store := NewMemoryStore()
	cat2, _ := NewCategory("tenant", "cat-2", "Second", 2)
	cat1, _ := NewCategory("tenant", "cat-1", "First", 1)
	_ = store.CreateCategory(cat2)
	_ = store.CreateCategory(cat1)

	categories := store.ListCategories("tenant")
	if len(categories) != 2 || categories[0].ID != "cat-1" || categories[1].ID != "cat-2" {
		t.Fatalf("unexpected order: %+v", categories)
	}
}
func TestMemoryStoreRejectsItemWithoutTenantCategory(t *testing.T) {
	store := NewMemoryStore()
	foreignCategory, _ := NewCategory("tenant-a", "cat-1", "Pizzas", 0)
	item := Item{TenantID: "tenant-b", ID: "item-1", CategoryID: foreignCategory.ID, Name: "Pizza", PriceCents: 1000}
	if err := store.CreateItem(item); err != ErrCategoryNotFound {
		t.Fatalf("err = %v, want %v", err, ErrCategoryNotFound)
	}
}

func TestMemoryStoreListsItemsByCategory(t *testing.T) {
	store := NewMemoryStore()
	catA, _ := NewCategory("tenant", "a", "A", 0)
	catB, _ := NewCategory("tenant", "b", "B", 1)
	_ = store.CreateCategory(catA)
	_ = store.CreateCategory(catB)
	item2, _ := NewItem("tenant", "i2", "Second", catA, 200, 2)
	item1, _ := NewItem("tenant", "i1", "First", catA, 100, 1)
	itemB, _ := NewItem("tenant", "ib", "Other", catB, 300, 0)
	_ = store.CreateItem(item2)
	_ = store.CreateItem(item1)
	_ = store.CreateItem(itemB)

	items := store.ListItems("tenant", "a")
	if len(items) != 2 || items[0].ID != "i1" || items[1].ID != "i2" {
		t.Fatalf("unexpected items: %+v", items)
	}
}
