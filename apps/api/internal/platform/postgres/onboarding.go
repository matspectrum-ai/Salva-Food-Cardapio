package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/tenant"
)

func (s *Store) Bootstrap(
	ctx context.Context,
	tenantValue tenant.Tenant,
	establishment tenant.Establishment,
	owner identity.Collaborator,
) error {
	tenantID, err := parseUUID(tenantValue.ID)
	if err != nil {
		return err
	}
	establishmentID, err := parseUUID(establishment.ID)
	if err != nil {
		return err
	}
	userID, err := parseUUID(owner.User.ID)
	if err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := s.q.WithTx(tx)

	if _, err = q.CreateTenant(ctx, db.CreateTenantParams{
		ID:   tenantID,
		Name: tenantValue.Name,
	}); err != nil {
		if pgConstraint(err, "23505", "tenants_pkey") {
			return onboarding.ErrTenantAlreadyExists
		}
		return err
	}

	if _, err = q.CreateEstablishment(ctx, db.CreateEstablishmentParams{
		ID:       establishmentID,
		TenantID: tenantID,
		Name:     establishment.Name,
		Timezone: establishment.Timezone,
	}); err != nil {
		if pgConstraint(err, "23505", "establishments_pkey") {
			return onboarding.ErrEstablishmentAlreadyExists
		}
		return err
	}
	if _, err = q.CreateAppUser(ctx, db.CreateAppUserParams{
		ID:           userID,
		Name:         owner.User.Name,
		Cpf:          owner.User.CPF,
		Email:        owner.User.Email,
		Phone:        owner.User.Phone,
		ImageUrl:     nullableText(owner.User.ImageURL),
		PasswordHash: owner.User.PasswordHash,
	}); err != nil {
		if pgConstraint(err, "23505", "") {
			return identity.ErrUserAlreadyExists
		}
		return err
	}

	if _, err = q.CreateTenantMembership(ctx, db.CreateTenantMembershipParams{
		TenantID:    tenantID,
		UserID:      userID,
		Title:       owner.Membership.Title,
		Status:      string(owner.Membership.Status),
		Permissions: permissionStrings(owner.Membership.Permissions),
	}); err != nil {
		if pgConstraint(err, "23505", "tenant_memberships_pkey") {
			return identity.ErrMembershipAlreadyExists
		}
		return err
	}

	return tx.Commit(ctx)
}

var _ onboarding.Repository = (*Store)(nil)
