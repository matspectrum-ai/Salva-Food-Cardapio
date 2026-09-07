package postgres

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
)

type Store struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	config.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{pool: pool, q: db.New(pool)}, nil
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, q: db.New(pool)}
}

func (s *Store) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func parseUUID(value string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid uuid %q: %w", value, err)
	}
	return id, nil
}

func formatUUID(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], id.Bytes[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], id.Bytes[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], id.Bytes[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], id.Bytes[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:36], id.Bytes[10:16])
	return string(buf)
}

func pgConstraint(err error, code, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != code {
		return false
	}
	return constraint == "" || pgErr.ConstraintName == constraint
}

func (s *Store) CreateCategory(ctx context.Context, category catalog.Category) error {
	id, err := parseUUID(category.ID)
	if err != nil {
		return err
	}
	tenantID, err := parseUUID(category.TenantID)
	if err != nil {
		return err
	}
	_, err = s.q.CreateCategory(ctx, db.CreateCategoryParams{
		ID: id, TenantID: tenantID, Name: category.Name,
		SortOrder: int32(category.SortOrder), SoldOut: category.SoldOut,
	})
	if pgConstraint(err, "23505", "") {
		return catalog.ErrAlreadyExists
	}
	return err
}

func (s *Store) GetCategory(ctx context.Context, tenantIDValue, idValue string) (catalog.Category, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return catalog.Category{}, false, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return catalog.Category{}, false, err
	}
	row, err := s.q.GetCategory(ctx, db.GetCategoryParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return catalog.Category{}, false, nil
	}
	if err != nil {
		return catalog.Category{}, false, err
	}
	return mapCategory(row), true, nil
}

func (s *Store) ListCategories(ctx context.Context, tenantIDValue string) ([]catalog.Category, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	rows, err := s.q.ListCategories(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	result := make([]catalog.Category, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapCategory(row))
	}
	return result, nil
}

func mapCategory(row db.CatalogCategory) catalog.Category {
	return catalog.Category{
		TenantID: formatUUID(row.TenantID), ID: formatUUID(row.ID),
		Name: row.Name, SortOrder: int(row.SortOrder), SoldOut: row.SoldOut,
	}
}

func (s *Store) CreateItem(ctx context.Context, item catalog.Item) error {
	id, err := parseUUID(item.ID)
	if err != nil {
		return err
	}
	tenantID, err := parseUUID(item.TenantID)
	if err != nil {
		return err
	}
	categoryID, err := parseUUID(item.CategoryID)
	if err != nil {
		return err
	}
	_, err = s.q.CreateItem(ctx, db.CreateItemParams{
		ID: id, TenantID: tenantID, CategoryID: categoryID,
		Name: item.Name, PriceCents: item.PriceCents,
		SortOrder: int32(item.SortOrder), SoldOut: item.SoldOut,
		SoldByWeight: item.SoldByWeight,
	})
	if pgConstraint(err, "23505", "") {
		return catalog.ErrAlreadyExists
	}
	if pgConstraint(err, "23503", "catalog_items_category_fk") {
		return catalog.ErrCategoryNotFound
	}
	return err
}

func (s *Store) GetItem(ctx context.Context, tenantIDValue, idValue string) (catalog.Item, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return catalog.Item{}, false, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return catalog.Item{}, false, err
	}
	row, err := s.q.GetItem(ctx, db.GetItemParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return catalog.Item{}, false, nil
	}
	if err != nil {
		return catalog.Item{}, false, err
	}
	return mapItem(row), true, nil
}

func (s *Store) ListItems(ctx context.Context, tenantIDValue, categoryIDValue string) ([]catalog.Item, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	var rows []db.CatalogItem
	if categoryIDValue == "" {
		rows, err = s.q.ListItems(ctx, tenantID)
	} else {
		categoryID, parseErr := parseUUID(categoryIDValue)
		if parseErr != nil {
			return nil, parseErr
		}
		rows, err = s.q.ListItemsByCategory(ctx, db.ListItemsByCategoryParams{
			TenantID: tenantID, CategoryID: categoryID,
		})
	}
	if err != nil {
		return nil, err
	}
	result := make([]catalog.Item, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapItem(row))
	}
	return result, nil
}

func mapItem(row db.CatalogItem) catalog.Item {
	return catalog.Item{
		TenantID: formatUUID(row.TenantID), ID: formatUUID(row.ID),
		CategoryID: formatUUID(row.CategoryID), Name: row.Name,
		PriceCents: row.PriceCents, SortOrder: int(row.SortOrder),
		SoldOut: row.SoldOut, SoldByWeight: row.SoldByWeight,
	}
}

func (s *Store) Replay(ctx context.Context, tenantIDValue, key, fingerprint string) (*orders.Order, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, false, err
	}
	record, err := s.q.GetOrderIdempotency(ctx, db.GetOrderIdempotencyParams{
		TenantID: tenantID, IdempotencyKey: key,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if record.RequestFingerprint != fingerprint {
		return nil, false, orders.ErrIdempotencyConflict
	}
	order, err := s.loadOrder(ctx, s.q, tenantID, record.OrderID)
	return order, err == nil, err
}

func (s *Store) Create(ctx context.Context, order *orders.Order, key, fingerprint string) (*orders.Order, bool, error) {
	tenantID, err := parseUUID(order.TenantID)
	if err != nil {
		return nil, false, err
	}
	orderID, err := parseUUID(order.ID)
	if err != nil {
		return nil, false, err
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := s.q.WithTx(tx)

	record, err := q.GetOrderIdempotency(ctx, db.GetOrderIdempotencyParams{
		TenantID: tenantID, IdempotencyKey: key,
	})
	if err == nil {
		if record.RequestFingerprint != fingerprint {
			return nil, false, orders.ErrIdempotencyConflict
		}
		existing, loadErr := s.loadOrder(ctx, q, tenantID, record.OrderID)
		if loadErr != nil {
			return nil, false, loadErr
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, false, err
		}
		return existing, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, false, err
	}

	_, err = q.CreateOrder(ctx, db.CreateOrderParams{
		ID: orderID, TenantID: tenantID, Source: order.Source,
		Status: string(order.Status), TotalCents: order.TotalCents,
	})
	if pgConstraint(err, "23505", "orders_pkey") {
		return nil, false, orders.ErrOrderAlreadyExists
	}
	if err != nil {
		return nil, false, err
	}

	for index, line := range order.Items {
		if index >= math.MaxInt32 {
			return nil, false, errors.New("too many order lines")
		}
		itemID, parseErr := parseUUID(line.ItemID)
		if parseErr != nil {
			return nil, false, parseErr
		}
		_, err = q.CreateOrderItem(ctx, db.CreateOrderItemParams{
			TenantID: tenantID, OrderID: orderID, LineNo: int32(index + 1),
			ItemID: itemID, NameSnapshot: line.Name,
			Quantity: line.Quantity, UnitPriceCents: line.UnitPriceCents,
		})
		if err != nil {
			return nil, false, err
		}
	}

	_, err = q.CreateOrderIdempotency(ctx, db.CreateOrderIdempotencyParams{
		TenantID: tenantID, IdempotencyKey: key,
		RequestFingerprint: fingerprint, OrderID: orderID,
	})
	if pgConstraint(err, "23505", "order_idempotency_pkey") {
		_ = tx.Rollback(ctx)
		existing, replay, replayErr := s.Replay(ctx, order.TenantID, key, fingerprint)
		if replayErr != nil {
			return nil, false, replayErr
		}
		if !replay {
			return nil, false, orders.ErrIdempotencyConflict
		}
		return existing, true, nil
	}
	if err != nil {
		return nil, false, err
	}

	created, err := s.loadOrder(ctx, q, tenantID, orderID)
	if err != nil {
		return nil, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, false, err
	}
	return created, false, nil
}

func (s *Store) Get(ctx context.Context, tenantIDValue, idValue string) (*orders.Order, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, false, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return nil, false, err
	}
	order, err := s.loadOrder(ctx, s.q, tenantID, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return order, true, nil
}

func (s *Store) List(ctx context.Context, tenantIDValue string) ([]*orders.Order, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	rows, err := s.q.ListOrders(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	result := make([]*orders.Order, 0, len(rows))
	for _, row := range rows {
		items, itemErr := s.q.ListOrderItems(ctx, db.ListOrderItemsParams{
			TenantID: tenantID, OrderID: row.ID,
		})
		if itemErr != nil {
			return nil, itemErr
		}
		result = append(result, mapOrder(row, items))
	}
	return result, nil
}

func (s *Store) Transition(ctx context.Context, tenantIDValue, idValue string, next orders.Status) (*orders.Order, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return nil, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := s.q.WithTx(tx)
	row, err := q.LockOrderForUpdate(ctx, db.LockOrderForUpdateParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, orders.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	items, err := q.ListOrderItems(ctx, db.ListOrderItemsParams{TenantID: tenantID, OrderID: id})
	if err != nil {
		return nil, err
	}
	order := mapOrder(row, items)
	if err := order.Transition(next); err != nil {
		return nil, err
	}
	_, err = q.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		TenantID: tenantID, ID: id, Status: string(order.Status),
	})
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *Store) loadOrder(ctx context.Context, q *db.Queries, tenantID, id pgtype.UUID) (*orders.Order, error) {
	row, err := q.GetOrder(ctx, db.GetOrderParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, err
	}
	items, err := q.ListOrderItems(ctx, db.ListOrderItemsParams{TenantID: tenantID, OrderID: id})
	if err != nil {
		return nil, err
	}
	return mapOrder(row, items), nil
}

func mapOrder(row db.Order, rows []db.OrderItem) *orders.Order {
	items := make([]orders.LineItemSnapshot, 0, len(rows))
	for _, line := range rows {
		items = append(items, orders.LineItemSnapshot{
			ItemID: formatUUID(line.ItemID), Name: line.NameSnapshot,
			Quantity: line.Quantity, UnitPriceCents: line.UnitPriceCents,
		})
	}
	return &orders.Order{
		TenantID:   formatUUID(row.TenantID),
		ID:         formatUUID(row.ID),
		Source:     row.Source,
		Status:     orders.Status(row.Status),
		Items:      items,
		TotalCents: row.TotalCents,
	}
}

var _ catalog.Repository = (*Store)(nil)
var _ orders.Repository = (*Store)(nil)

func permissionStrings(input []identity.Permission) []string {
	result := make([]string, len(input))
	for i, permission := range input {
		result[i] = string(permission)
	}
	return result
}

func mapIdentityUser(row db.AppUser) (identity.User, error) {
	user, err := identity.NewUser(
		formatUUID(row.ID), row.Name, row.Cpf, row.Email, row.Phone, row.PasswordHash,
	)
	if err != nil {
		return identity.User{}, err
	}
	if row.ImageUrl.Valid {
		user.ImageURL = row.ImageUrl.String
	}
	return user, nil
}

func mapIdentityMembership(row db.TenantMembership) (identity.Membership, error) {
	permissions := make([]identity.Permission, len(row.Permissions))
	for i, permission := range row.Permissions {
		permissions[i] = identity.Permission(permission)
	}
	membership, err := identity.NewMembership(
		formatUUID(row.TenantID), formatUUID(row.UserID), row.Title, permissions,
	)
	if err != nil {
		return identity.Membership{}, err
	}
	if err := membership.SetStatus(identity.MembershipStatus(row.Status)); err != nil {
		return identity.Membership{}, err
	}
	return membership, nil
}

func nullableText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func (s *Store) CreateCollaborator(ctx context.Context, collaborator identity.Collaborator) error {
	userID, err := parseUUID(collaborator.User.ID)
	if err != nil {
		return err
	}
	tenantID, err := parseUUID(collaborator.Membership.TenantID)
	if err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := s.q.WithTx(tx)

	_, err = q.CreateAppUser(ctx, db.CreateAppUserParams{
		ID: userID, Name: collaborator.User.Name, Cpf: collaborator.User.CPF,
		Email: collaborator.User.Email, Phone: collaborator.User.Phone,
		ImageUrl:     nullableText(collaborator.User.ImageURL),
		PasswordHash: collaborator.User.PasswordHash,
	})
	if pgConstraint(err, "23505", "") {
		return identity.ErrUserAlreadyExists
	}
	if err != nil {
		return err
	}
	_, err = q.CreateTenantMembership(ctx, db.CreateTenantMembershipParams{
		TenantID: tenantID, UserID: userID, Title: collaborator.Membership.Title,
		Status:      string(collaborator.Membership.Status),
		Permissions: permissionStrings(collaborator.Membership.Permissions),
	})
	if pgConstraint(err, "23505", "tenant_memberships_pkey") {
		return identity.ErrMembershipAlreadyExists
	}
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (identity.User, bool, error) {
	row, err := s.q.GetAppUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return identity.User{}, false, nil
	}
	if err != nil {
		return identity.User{}, false, err
	}
	user, err := mapIdentityUser(row)
	if err != nil {
		return identity.User{}, false, err
	}
	return user, true, nil
}

func (s *Store) GetMembership(ctx context.Context, tenantIDValue, userIDValue string) (identity.Membership, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return identity.Membership{}, false, err
	}
	userID, err := parseUUID(userIDValue)
	if err != nil {
		return identity.Membership{}, false, err
	}
	row, err := s.q.GetTenantMembership(ctx, db.GetTenantMembershipParams{TenantID: tenantID, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return identity.Membership{}, false, nil
	}
	if err != nil {
		return identity.Membership{}, false, err
	}
	membership, err := mapIdentityMembership(row)
	if err != nil {
		return identity.Membership{}, false, err
	}
	return membership, true, nil
}

func (s *Store) ListCollaborators(ctx context.Context, tenantIDValue, search string) ([]identity.Collaborator, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	rows, err := s.q.ListCollaborators(ctx, db.ListCollaboratorsParams{TenantID: tenantID, Search: search})
	if err != nil {
		return nil, err
	}
	result := make([]identity.Collaborator, 0, len(rows))
	for _, row := range rows {
		user, err := mapIdentityUser(db.AppUser{
			ID: row.ID, Name: row.Name, Cpf: row.Cpf, Email: row.Email,
			Phone: row.Phone, ImageUrl: row.ImageUrl, PasswordHash: row.PasswordHash,
		})
		if err != nil {
			return nil, err
		}
		membership, err := mapIdentityMembership(db.TenantMembership{
			TenantID: row.TenantID, UserID: row.ID, Title: row.Title,
			Status: row.Status, Permissions: row.Permissions,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, identity.Collaborator{User: user, Membership: membership})
	}
	return result, nil
}

func (s *Store) SetMembershipStatus(ctx context.Context, tenantIDValue, userIDValue string, status identity.MembershipStatus) error {
	if status != identity.MembershipActive && status != identity.MembershipInactive {
		return identity.ErrInvalidStatus
	}
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return err
	}
	userID, err := parseUUID(userIDValue)
	if err != nil {
		return err
	}
	_, err = s.q.SetTenantMembershipStatus(ctx, db.SetTenantMembershipStatusParams{
		TenantID: tenantID, UserID: userID, Status: string(status),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return identity.ErrCollaboratorNotFound
	}
	return err
}

var _ identity.Repository = (*Store)(nil)
