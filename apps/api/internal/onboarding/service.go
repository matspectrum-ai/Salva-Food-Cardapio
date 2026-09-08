package onboarding

import (
	"context"
	"errors"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/tenant"
)

const OwnerTitle = "Proprietário"

var (
	ErrTenantAlreadyExists        = errors.New("tenant already exists")
	ErrEstablishmentAlreadyExists = errors.New("establishment already exists")
)

type OwnerInput struct {
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Input struct {
	TenantName        string     `json:"tenant_name"`
	EstablishmentName string     `json:"establishment_name"`
	Timezone          string     `json:"timezone"`
	Owner             OwnerInput `json:"owner"`
}

type Result struct {
	Tenant        tenant.Tenant         `json:"tenant"`
	Establishment tenant.Establishment  `json:"establishment"`
	Owner         identity.Collaborator `json:"owner"`
}

type Repository interface {
	Bootstrap(context.Context, tenant.Tenant, tenant.Establishment, identity.Collaborator) error
}

type Service struct {
	repo   Repository
	hasher identity.PasswordHasher
	newID  func() string
}

func NewService(repo Repository, hasher identity.PasswordHasher, newID func() string) *Service {
	return &Service{repo: repo, hasher: hasher, newID: newID}
}

func (s *Service) Create(ctx context.Context, input Input) (Result, error) {
	if err := identity.ValidatePassword(input.Owner.Password); err != nil {
		return Result{}, err
	}
	tenantValue, err := tenant.New(s.newID(), input.TenantName)
	if err != nil {
		return Result{}, err
	}
	establishment, err := tenant.NewEstablishment(
		s.newID(), tenantValue.ID, input.EstablishmentName, input.Timezone,
	)
	if err != nil {
		return Result{}, err
	}
	hash, err := s.hasher.Hash(input.Owner.Password)
	if err != nil {
		return Result{}, err
	}
	user, err := identity.NewUser(
		s.newID(), input.Owner.Name, input.Owner.CPF,
		input.Owner.Email, input.Owner.Phone, hash,
	)
	if err != nil {
		return Result{}, err
	}
	membership, err := identity.NewMembership(
		tenantValue.ID, user.ID, OwnerTitle, identity.AllPermissions(),
	)
	if err != nil {
		return Result{}, err
	}
	owner := identity.Collaborator{User: user, Membership: membership}
	if err := s.repo.Bootstrap(ctx, tenantValue, establishment, owner); err != nil {
		return Result{}, err
	}
	return Result{Tenant: tenantValue, Establishment: establishment, Owner: owner}, nil
}
