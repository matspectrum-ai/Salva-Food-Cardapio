package orders

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/customers"
)

var (
	ErrIdempotencyKeyRequired  = errors.New("idempotency key is required")
	ErrCatalogItemUnavailable  = errors.New("catalog item is unavailable")
	ErrCustomerNotFound        = errors.New("customer not found")
	ErrAddressNotFound         = errors.New("address not found")
	ErrAddressCustomerMismatch = errors.New("address does not belong to customer")
	ErrDeliveryAddressRequired = errors.New("delivery address is required")
)

type CatalogProduct struct {
	ID              string
	Name            string
	PriceCents      int64
	SoldOut         bool
	CategorySoldOut bool
}

type CatalogReader interface {
	LookupOrderProduct(context.Context, string, string) (CatalogProduct, error)
}

type IDGenerator func() string

type CreateLineInput struct {
	ItemID   string `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

type CreateInput struct {
	Source      string            `json:"source"`
	Items       []CreateLineInput `json:"items"`
	CustomerID  string            `json:"customer_id,omitempty"`
	AddressID   string            `json:"address_id,omitempty"`
	Fulfillment FulfillmentType   `json:"fulfillment_type,omitempty"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	Notes       string            `json:"notes,omitempty"`
}

type CustomerReader interface {
	GetCustomer(context.Context, string, string) (customers.Customer, bool, error)
	GetAddress(context.Context, string, string) (customers.Address, bool, error)
}

type Service struct {
	store    Repository
	catalog  CatalogReader
	customer CustomerReader
	newID    IDGenerator
}

func NewService(store Repository, catalog CatalogReader, newID IDGenerator, readers ...CustomerReader) *Service {
	var customer CustomerReader
	if len(readers) > 0 {
		customer = readers[0]
	}
	return &Service{store: store, catalog: catalog, customer: customer, newID: newID}
}
func (s *Service) Create(ctx context.Context, tenantID, idempotencyKey string, input CreateInput) (*Order, bool, error) {
	if strings.TrimSpace(idempotencyKey) == "" {
		return nil, false, ErrIdempotencyKeyRequired
	}
	fingerprint, err := fingerprintInput(input)
	if err != nil {
		return nil, false, err
	}
	if existing, replay, err := s.store.Replay(ctx, tenantID, idempotencyKey, fingerprint); err != nil || replay {
		return existing, replay, err
	}

	items := make([]LineItemSnapshot, 0, len(input.Items))
	for _, requested := range input.Items {
		product, err := s.catalog.LookupOrderProduct(ctx, tenantID, requested.ItemID)
		if err != nil {
			return nil, false, err
		}
		if product.SoldOut || product.CategorySoldOut {
			return nil, false, ErrCatalogItemUnavailable
		}
		items = append(items, LineItemSnapshot{
			ItemID:         product.ID,
			Name:           product.Name,
			Quantity:       requested.Quantity,
			UnitPriceCents: product.PriceCents,
		})
	}

	order, err := NewOrder(tenantID, s.newID(), input.Source, items)
	if err != nil {
		return nil, false, err
	}
	if input.Fulfillment != "" {
		if err := order.SetFulfillment(input.Fulfillment); err != nil {
			return nil, false, err
		}
	}
	order.ScheduledAt = input.ScheduledAt
	order.Notes = strings.TrimSpace(input.Notes)
	if input.CustomerID != "" {
		if s.customer == nil {
			return nil, false, ErrCustomerNotFound
		}
		c, found, err := s.customer.GetCustomer(ctx, tenantID, input.CustomerID)
		if err != nil {
			return nil, false, err
		}
		if !found {
			return nil, false, ErrCustomerNotFound
		}
		order.Customer = &CustomerSnapshot{ID: c.ID, Name: c.Name, Phone: c.Phone, Email: c.Email}
	}
	if input.AddressID != "" {
		if s.customer == nil {
			return nil, false, ErrAddressNotFound
		}
		a, found, err := s.customer.GetAddress(ctx, tenantID, input.AddressID)
		if err != nil {
			return nil, false, err
		}
		if !found {
			return nil, false, ErrAddressNotFound
		}
		if order.Customer == nil || a.CustomerID != order.Customer.ID {
			return nil, false, ErrAddressCustomerMismatch
		}
		order.Address = &AddressSnapshot{ID: a.ID, Label: a.Label, Street: a.Street, Number: a.Number, Complement: a.Complement, Neighborhood: a.Neighborhood, City: a.City, State: a.State, PostalCode: a.PostalCode, Reference: a.Reference}
	}
	if order.Fulfillment == FulfillmentDelivery && order.Address == nil {
		return nil, false, ErrDeliveryAddressRequired
	}
	return s.store.Create(ctx, order, idempotencyKey, fingerprint)
}

func fingerprintInput(input CreateInput) (string, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
