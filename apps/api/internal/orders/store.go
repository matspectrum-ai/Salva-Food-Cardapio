package orders

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyExists  = errors.New("order already exists")
	ErrIdempotencyConflict = errors.New("idempotency key reused with different payload")
)

type idempotencyRecord struct {
	OrderID     string
	Fingerprint string
}

type MemoryStore struct {
	mu          sync.RWMutex
	orders      map[string]map[string]*Order
	idempotency map[string]map[string]idempotencyRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		orders:      make(map[string]map[string]*Order),
		idempotency: make(map[string]map[string]idempotencyRecord),
	}
}

func (s *MemoryStore) Create(order *Order, key, fingerprint string) (*Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.orders[order.TenantID] == nil {
		s.orders[order.TenantID] = make(map[string]*Order)
	}
	if s.idempotency[order.TenantID] == nil {
		s.idempotency[order.TenantID] = make(map[string]idempotencyRecord)
	}
	if record, exists := s.idempotency[order.TenantID][key]; exists {
		if record.Fingerprint != fingerprint {
			return nil, false, ErrIdempotencyConflict
		}
		existing := s.orders[order.TenantID][record.OrderID]
		return cloneOrder(existing), true, nil
	}
	if _, exists := s.orders[order.TenantID][order.ID]; exists {
		return nil, false, ErrOrderAlreadyExists
	}
	stored := cloneOrder(order)
	s.orders[order.TenantID][order.ID] = stored
	s.idempotency[order.TenantID][key] = idempotencyRecord{OrderID: order.ID, Fingerprint: fingerprint}
	return cloneOrder(stored), false, nil
}

func (s *MemoryStore) Get(tenantID, id string) (*Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[tenantID][id]
	if !ok {
		return nil, false
	}
	return cloneOrder(order), true
}

func (s *MemoryStore) List(tenantID string) []*Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Order, 0, len(s.orders[tenantID]))
	for _, order := range s.orders[tenantID] {
		result = append(result, cloneOrder(order))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (s *MemoryStore) Transition(tenantID, id string, next Status) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[tenantID][id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if err := order.Transition(next); err != nil {
		return nil, err
	}
	return cloneOrder(order), nil
}

func cloneOrder(order *Order) *Order {
	if order == nil {
		return nil
	}
	clone := *order
	clone.Items = make([]LineItemSnapshot, len(order.Items))
	copy(clone.Items, order.Items)
	return &clone
}
func (s *MemoryStore) Replay(tenantID, key, fingerprint string) (*Order, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, exists := s.idempotency[tenantID][key]
	if !exists {
		return nil, false, nil
	}
	if record.Fingerprint != fingerprint {
		return nil, false, ErrIdempotencyConflict
	}
	order := s.orders[tenantID][record.OrderID]
	return cloneOrder(order), true, nil
}
