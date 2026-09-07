package catalog

import (
	"errors"
	"strings"
)

var (
	ErrTenantRequired    = errors.New("tenant id is required")
	ErrIDRequired        = errors.New("id is required")
	ErrNameRequired      = errors.New("name is required")
	ErrNegativeSortOrder = errors.New("sort order cannot be negative")
	ErrNegativePrice     = errors.New("price cannot be negative")
	ErrTenantMismatch    = errors.New("category belongs to another tenant")
)

type Category struct {
	TenantID  string `json:"tenant_id"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	SoldOut   bool   `json:"sold_out"`
}

func NewCategory(tenantID, id, name string, sortOrder int) (Category, error) {
	if strings.TrimSpace(tenantID) == "" {
		return Category{}, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return Category{}, ErrIDRequired
	}
	if strings.TrimSpace(name) == "" {
		return Category{}, ErrNameRequired
	}
	if sortOrder < 0 {
		return Category{}, ErrNegativeSortOrder
	}
	return Category{TenantID: tenantID, ID: id, Name: strings.TrimSpace(name), SortOrder: sortOrder}, nil
}
func (c *Category) SetSortOrder(sortOrder int) error {
	if sortOrder < 0 {
		return ErrNegativeSortOrder
	}
	c.SortOrder = sortOrder
	return nil
}

func (c *Category) SetSoldOut(soldOut bool) {
	c.SoldOut = soldOut
}

type Item struct {
	TenantID     string `json:"tenant_id"`
	ID           string `json:"id"`
	CategoryID   string `json:"category_id"`
	Name         string `json:"name"`
	PriceCents   int64  `json:"price_cents"`
	SortOrder    int    `json:"sort_order"`
	SoldOut      bool   `json:"sold_out"`
	SoldByWeight bool   `json:"sold_by_weight"`
}

func NewItem(tenantID, id, name string, category Category, priceCents int64, sortOrder int) (Item, error) {
	if strings.TrimSpace(tenantID) == "" {
		return Item{}, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return Item{}, ErrIDRequired
	}
	if strings.TrimSpace(name) == "" {
		return Item{}, ErrNameRequired
	}
	if tenantID != category.TenantID {
		return Item{}, ErrTenantMismatch
	}
	if priceCents < 0 {
		return Item{}, ErrNegativePrice
	}
	if sortOrder < 0 {
		return Item{}, ErrNegativeSortOrder
	}
	return Item{
		TenantID:   tenantID,
		ID:         id,
		CategoryID: category.ID,
		Name:       strings.TrimSpace(name),
		PriceCents: priceCents,
		SortOrder:  sortOrder,
	}, nil
}

func (i *Item) SetPrice(priceCents int64) error {
	if priceCents < 0 {
		return ErrNegativePrice
	}
	i.PriceCents = priceCents
	return nil
}

func (i *Item) SetSortOrder(sortOrder int) error {
	if sortOrder < 0 {
		return ErrNegativeSortOrder
	}
	i.SortOrder = sortOrder
	return nil
}

func (i *Item) SetSoldOut(soldOut bool) {
	i.SoldOut = soldOut
}
