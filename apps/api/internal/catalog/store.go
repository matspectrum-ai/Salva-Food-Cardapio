package catalog

import (
	"errors"
	"sort"
	"sync"
)

var (
	ErrAlreadyExists    = errors.New("catalog entity already exists")
	ErrCategoryNotFound = errors.New("category not found")
	ErrItemNotFound     = errors.New("item not found")
)

type MemoryStore struct {
	mu         sync.RWMutex
	categories map[string]map[string]Category
	items      map[string]map[string]Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		categories: make(map[string]map[string]Category),
		items:      make(map[string]map[string]Item),
	}
}

func (s *MemoryStore) CreateCategory(category Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.categories[category.TenantID] == nil {
		s.categories[category.TenantID] = make(map[string]Category)
	}
	if _, exists := s.categories[category.TenantID][category.ID]; exists {
		return ErrAlreadyExists
	}
	s.categories[category.TenantID][category.ID] = category
	return nil
}
func (s *MemoryStore) GetCategory(tenantID, id string) (Category, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	category, ok := s.categories[tenantID][id]
	return category, ok
}

func (s *MemoryStore) ListCategories(tenantID string) []Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Category, 0, len(s.categories[tenantID]))
	for _, category := range s.categories[tenantID] {
		result = append(result, category)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SortOrder == result[j].SortOrder {
			return result[i].ID < result[j].ID
		}
		return result[i].SortOrder < result[j].SortOrder
	})
	return result
}

func (s *MemoryStore) CreateItem(item Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[item.TenantID][item.CategoryID]; !ok {
		return ErrCategoryNotFound
	}
	if s.items[item.TenantID] == nil {
		s.items[item.TenantID] = make(map[string]Item)
	}
	if _, exists := s.items[item.TenantID][item.ID]; exists {
		return ErrAlreadyExists
	}
	s.items[item.TenantID][item.ID] = item
	return nil
}
func (s *MemoryStore) GetItem(tenantID, id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[tenantID][id]
	return item, ok
}

func (s *MemoryStore) ListItems(tenantID, categoryID string) []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Item, 0)
	for _, item := range s.items[tenantID] {
		if categoryID == "" || item.CategoryID == categoryID {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SortOrder == result[j].SortOrder {
			return result[i].ID < result[j].ID
		}
		return result[i].SortOrder < result[j].SortOrder
	})
	return result
}
