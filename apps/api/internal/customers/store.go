package customers

import (
	"context"
	"errors"
	"sort"
	"sync"
)

var (
	ErrAlreadyExists    = errors.New("customer entity already exists")
	ErrCustomerNotFound = errors.New("customer not found")
)

type Repository interface {
	CreateCustomer(context.Context, Customer) error
	GetCustomer(context.Context, string, string) (Customer, bool, error)
	ListCustomers(context.Context, string) ([]Customer, error)
	CreateAddress(context.Context, Address) error
	GetAddress(context.Context, string, string) (Address, bool, error)
	ListAddresses(context.Context, string, string) ([]Address, error)
}

type MemoryStore struct {
	mu        sync.RWMutex
	customers map[string]map[string]Customer
	addresses map[string]map[string]Address
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{customers: make(map[string]map[string]Customer), addresses: make(map[string]map[string]Address)}
}

func (s *MemoryStore) CreateCustomer(_ context.Context, c Customer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.customers[c.TenantID] == nil {
		s.customers[c.TenantID] = make(map[string]Customer)
	}
	if _, ok := s.customers[c.TenantID][c.ID]; ok {
		return ErrAlreadyExists
	}
	s.customers[c.TenantID][c.ID] = c
	return nil
}

func (s *MemoryStore) GetCustomer(_ context.Context, tenantID, id string) (Customer, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.customers[tenantID][id]
	return c, ok, nil
}

func (s *MemoryStore) ListCustomers(_ context.Context, tenantID string) ([]Customer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Customer, 0, len(s.customers[tenantID]))
	for _, c := range s.customers[tenantID] {
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *MemoryStore) CreateAddress(_ context.Context, a Address) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.customers[a.TenantID][a.CustomerID]; !ok {
		return ErrCustomerNotFound
	}
	if s.addresses[a.TenantID] == nil {
		s.addresses[a.TenantID] = make(map[string]Address)
	}
	if _, ok := s.addresses[a.TenantID][a.ID]; ok {
		return ErrAlreadyExists
	}
	s.addresses[a.TenantID][a.ID] = a
	return nil
}

func (s *MemoryStore) GetAddress(_ context.Context, tenantID, id string) (Address, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.addresses[tenantID][id]
	return a, ok, nil
}

func (s *MemoryStore) ListAddresses(_ context.Context, tenantID, customerID string) ([]Address, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Address, 0)
	for _, a := range s.addresses[tenantID] {
		if customerID == "" || a.CustomerID == customerID {
			result = append(result, a)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

var _ Repository = (*MemoryStore)(nil)
