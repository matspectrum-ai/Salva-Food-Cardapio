package catalog

import "testing"

func TestNewCategoryValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name, tenantID, id, categoryName string
		sortOrder                        int
	}{
		{"missing tenant", "", "cat-1", "Pizzas", 0},
		{"missing id", "tenant-1", "", "Pizzas", 0},
		{"missing name", "tenant-1", "cat-1", " ", 0},
		{"negative order", "tenant-1", "cat-1", "Pizzas", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewCategory(tt.tenantID, tt.id, tt.categoryName, tt.sortOrder); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestCategoryCanBeReorderedAndSoldOut(t *testing.T) {
	category, err := NewCategory("tenant-1", "cat-1", "Pizzas", 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := category.SetSortOrder(0); err != nil {
		t.Fatal(err)
	}
	category.SetSoldOut(true)

	if category.SortOrder != 0 || !category.SoldOut {
		t.Fatalf("unexpected category state: %+v", category)
	}
}
func TestNewItemEnforcesTenantAndMoneyInvariants(t *testing.T) {
	category, err := NewCategory("tenant-1", "cat-1", "Pizzas", 0)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := NewItem("tenant-2", "item-1", "Calabresa", category, 3590, 0); err == nil {
		t.Fatal("expected cross-tenant category error")
	}
	if _, err := NewItem("tenant-1", "item-1", "Calabresa", category, -1, 0); err == nil {
		t.Fatal("expected negative price error")
	}
	if _, err := NewItem("tenant-1", "item-1", " ", category, 3590, 0); err == nil {
		t.Fatal("expected blank name error")
	}
}

func TestItemPriceOrderAndAvailabilityAreIndependent(t *testing.T) {
	category, _ := NewCategory("tenant-1", "cat-1", "Pizzas", 0)
	item, err := NewItem("tenant-1", "item-1", "Calabresa", category, 3590, 4)
	if err != nil {
		t.Fatal(err)
	}

	if err := item.SetPrice(3990); err != nil {
		t.Fatal(err)
	}
	if err := item.SetSortOrder(1); err != nil {
		t.Fatal(err)
	}
	item.SetSoldOut(true)

	if item.PriceCents != 3990 || item.SortOrder != 1 || !item.SoldOut {
		t.Fatalf("unexpected item state: %+v", item)
	}
	if err := item.SetPrice(-1); err == nil {
		t.Fatal("expected negative price error")
	}
}
