package onboarding

import (
	"context"
	"sync"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/tenant"
)

type MemoryRepository struct {
	mu             sync.RWMutex
	identityRepo   identity.Repository
	tenants        map[string]tenant.Tenant
	establishments map[string]map[string]tenant.Establishment
}

func NewMemoryRepository(identityRepo identity.Repository) *MemoryRepository {
	return &MemoryRepository{
		identityRepo:   identityRepo,
		tenants:        make(map[string]tenant.Tenant),
		establishments: make(map[string]map[string]tenant.Establishment),
	}
}

func (s *MemoryRepository) Bootstrap(
	ctx context.Context,
	tenantValue tenant.Tenant,
	establishment tenant.Establishment,
	owner identity.Collaborator,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tenants[tenantValue.ID]; exists {
		return ErrTenantAlreadyExists
	}
	if _, found, err := s.identityRepo.FindUserByEmail(ctx, owner.User.Email); err != nil {
		return err
	} else if found {
		return identity.ErrUserAlreadyExists
	}

	if s.establishments[tenantValue.ID] == nil {
		s.establishments[tenantValue.ID] = make(map[string]tenant.Establishment)
	}
	if _, exists := s.establishments[tenantValue.ID][establishment.ID]; exists {
		return ErrEstablishmentAlreadyExists
	}

	s.tenants[tenantValue.ID] = tenantValue
	s.establishments[tenantValue.ID][establishment.ID] = establishment
	if err := s.identityRepo.CreateCollaborator(ctx, owner); err != nil {
		delete(s.establishments[tenantValue.ID], establishment.ID)
		delete(s.establishments, tenantValue.ID)
		delete(s.tenants, tenantValue.ID)
		return err
	}
	return nil
}

func (s *MemoryRepository) GetTenant(id string) (tenant.Tenant, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.tenants[id]
	return value, ok
}

func (s *MemoryRepository) GetEstablishment(tenantID, id string) (tenant.Establishment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.establishments[tenantID][id]
	return value, ok
}

var _ Repository = (*MemoryRepository)(nil)
