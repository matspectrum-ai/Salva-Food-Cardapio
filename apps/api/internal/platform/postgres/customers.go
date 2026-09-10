package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/customers"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
)

func (s *Store) CreateCustomer(ctx context.Context, c customers.Customer) error {
	tenantID, err := parseUUID(c.TenantID)
	if err != nil {
		return err
	}
	id, err := parseUUID(c.ID)
	if err != nil {
		return err
	}
	_, err = s.q.CreateCustomer(ctx, db.CreateCustomerParams{TenantID: tenantID, ID: id, Name: c.Name, Phone: c.Phone, Email: nullableText(c.Email)})
	if pgConstraint(err, "23505", "") {
		return customers.ErrAlreadyExists
	}
	return err
}

func (s *Store) GetCustomer(ctx context.Context, tenantIDValue, idValue string) (customers.Customer, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return customers.Customer{}, false, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return customers.Customer{}, false, err
	}
	row, err := s.q.GetCustomer(ctx, db.GetCustomerParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return customers.Customer{}, false, nil
	}
	if err != nil {
		return customers.Customer{}, false, err
	}
	return mapCustomer(row), true, nil
}

func (s *Store) ListCustomers(ctx context.Context, tenantIDValue string) ([]customers.Customer, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	rows, err := s.q.ListCustomers(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	result := make([]customers.Customer, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapCustomer(row))
	}
	return result, nil
}

func mapCustomer(row db.Customer) customers.Customer {
	c := customers.Customer{TenantID: formatUUID(row.TenantID), ID: formatUUID(row.ID), Name: row.Name, Phone: row.Phone}
	if row.Email.Valid {
		c.Email = row.Email.String
	}
	return c
}

func (s *Store) CreateAddress(ctx context.Context, a customers.Address) error {
	tenantID, err := parseUUID(a.TenantID)
	if err != nil {
		return err
	}
	id, err := parseUUID(a.ID)
	if err != nil {
		return err
	}
	customerID, err := parseUUID(a.CustomerID)
	if err != nil {
		return err
	}
	_, err = s.q.CreateCustomerAddress(ctx, db.CreateCustomerAddressParams{
		TenantID: tenantID, ID: id, CustomerID: customerID, Label: nullableText(a.Label),
		Street: a.Street, Number: a.Number, Complement: nullableText(a.Complement),
		Neighborhood: a.Neighborhood, City: a.City, State: a.State,
		PostalCode: a.PostalCode, Reference: nullableText(a.Reference),
	})
	if pgConstraint(err, "23505", "") {
		return customers.ErrAlreadyExists
	}
	if pgConstraint(err, "23503", "customer_addresses_customer_fk") {
		return customers.ErrCustomerNotFound
	}
	return err
}

func (s *Store) GetAddress(ctx context.Context, tenantIDValue, idValue string) (customers.Address, bool, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return customers.Address{}, false, err
	}
	id, err := parseUUID(idValue)
	if err != nil {
		return customers.Address{}, false, err
	}
	row, err := s.q.GetCustomerAddress(ctx, db.GetCustomerAddressParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return customers.Address{}, false, nil
	}
	if err != nil {
		return customers.Address{}, false, err
	}
	return mapAddress(row), true, nil
}

func (s *Store) ListAddresses(ctx context.Context, tenantIDValue, customerIDValue string) ([]customers.Address, error) {
	tenantID, err := parseUUID(tenantIDValue)
	if err != nil {
		return nil, err
	}
	customerID, err := parseUUID(customerIDValue)
	if err != nil {
		return nil, err
	}
	rows, err := s.q.ListCustomerAddresses(ctx, db.ListCustomerAddressesParams{TenantID: tenantID, CustomerID: customerID})
	if err != nil {
		return nil, err
	}
	result := make([]customers.Address, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapAddress(row))
	}
	return result, nil
}
func mapAddress(row db.CustomerAddress) customers.Address {
	a := customers.Address{TenantID: formatUUID(row.TenantID), ID: formatUUID(row.ID), CustomerID: formatUUID(row.CustomerID), Street: row.Street, Number: row.Number, Neighborhood: row.Neighborhood, City: row.City, State: row.State, PostalCode: row.PostalCode}
	if row.Label.Valid {
		a.Label = row.Label.String
	}
	if row.Complement.Valid {
		a.Complement = row.Complement.String
	}
	if row.Reference.Valid {
		a.Reference = row.Reference.String
	}
	return a
}

var _ customers.Repository = (*Store)(nil)
