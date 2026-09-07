package orders

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrIdempotencyKeyRequired = errors.New("idempotency key is required")
	ErrCatalogItemUnavailable = errors.New("catalog item is unavailable")
)

type CatalogProduct struct {
	ID              string
	Name            string
	PriceCents      int64
	SoldOut         bool
	CategorySoldOut bool
}

type CatalogReader interface {
	LookupOrderProduct(tenantID, itemID string) (CatalogProduct, error)
}

type IDGenerator func() string

type CreateLineInput struct {
	ItemID   string `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

type CreateInput struct {
	Source string            `json:"source"`
	Items  []CreateLineInput `json:"items"`
}

type Service struct {
	store   *MemoryStore
	catalog CatalogReader
	newID   IDGenerator
}

func NewService(store *MemoryStore, catalog CatalogReader, newID IDGenerator) *Service {
	return &Service{store: store, catalog: catalog, newID: newID}
}
func (s *Service) Create(tenantID, idempotencyKey string, input CreateInput) (*Order, bool, error) {
	if strings.TrimSpace(idempotencyKey) == "" {
		return nil, false, ErrIdempotencyKeyRequired
	}
	fingerprint, err := fingerprintInput(input)
	if err != nil {
		return nil, false, err
	}
	if existing, replay, err := s.store.Replay(tenantID, idempotencyKey, fingerprint); err != nil || replay {
		return existing, replay, err
	}

	items := make([]LineItemSnapshot, 0, len(input.Items))
	for _, requested := range input.Items {
		product, err := s.catalog.LookupOrderProduct(tenantID, requested.ItemID)
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
	return s.store.Create(order, idempotencyKey, fingerprint)
}

func fingerprintInput(input CreateInput) (string, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
