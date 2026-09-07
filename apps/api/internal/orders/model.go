package orders

import (
	"errors"
	"math"
	"strings"
)

var (
	ErrTenantRequired    = errors.New("tenant id is required")
	ErrIDRequired        = errors.New("id is required")
	ErrSourceRequired    = errors.New("order source is required")
	ErrItemsRequired     = errors.New("at least one item is required")
	ErrInvalidLineItem   = errors.New("invalid line item")
	ErrTotalOverflow     = errors.New("order total overflow")
	ErrInvalidStatus     = errors.New("invalid order status")
	ErrInvalidTransition = errors.New("invalid order transition")
)

type Status string

const (
	StatusAnalysis   Status = "ANALYSIS"
	StatusProduction Status = "PRODUCTION"
	StatusReady      Status = "READY"
	StatusFinalized  Status = "FINALIZED"
	StatusCancelled  Status = "CANCELLED"
)

type LineItemSnapshot struct {
	ItemID         string `json:"item_id"`
	Name           string `json:"name"`
	Quantity       int64  `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

type Order struct {
	TenantID   string             `json:"tenant_id"`
	ID         string             `json:"id"`
	Source     string             `json:"source"`
	Status     Status             `json:"status"`
	Items      []LineItemSnapshot `json:"items"`
	TotalCents int64              `json:"total_cents"`
}

func NewOrder(tenantID, id, source string, items []LineItemSnapshot) (*Order, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return nil, ErrIDRequired
	}
	if strings.TrimSpace(source) == "" {
		return nil, ErrSourceRequired
	}
	if len(items) == 0 {
		return nil, ErrItemsRequired
	}

	cloned := make([]LineItemSnapshot, len(items))
	copy(cloned, items)
	var total int64
	for _, item := range cloned {
		if strings.TrimSpace(item.ItemID) == "" || strings.TrimSpace(item.Name) == "" || item.Quantity <= 0 || item.UnitPriceCents < 0 {
			return nil, ErrInvalidLineItem
		}
		if item.UnitPriceCents != 0 && item.Quantity > math.MaxInt64/item.UnitPriceCents {
			return nil, ErrTotalOverflow
		}
		lineTotal := item.Quantity * item.UnitPriceCents
		if total > math.MaxInt64-lineTotal {
			return nil, ErrTotalOverflow
		}
		total += lineTotal
	}

	return &Order{
		TenantID:   tenantID,
		ID:         id,
		Source:     strings.TrimSpace(source),
		Status:     StatusAnalysis,
		Items:      cloned,
		TotalCents: total,
	}, nil
}
func (o *Order) Transition(next Status) error {
	if !validStatus(next) {
		return ErrInvalidStatus
	}
	if next == o.Status {
		return nil
	}
	if o.Status == StatusFinalized || o.Status == StatusCancelled {
		return ErrInvalidTransition
	}
	if next == StatusCancelled {
		o.Status = next
		return nil
	}

	allowed := map[Status]Status{
		StatusAnalysis:   StatusProduction,
		StatusProduction: StatusReady,
		StatusReady:      StatusFinalized,
	}
	if allowed[o.Status] != next {
		return ErrInvalidTransition
	}
	o.Status = next
	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusAnalysis, StatusProduction, StatusReady, StatusFinalized, StatusCancelled:
		return true
	default:
		return false
	}
}
