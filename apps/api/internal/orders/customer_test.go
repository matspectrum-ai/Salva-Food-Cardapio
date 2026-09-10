package orders

import (
	"context"
	"testing"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/customers"
)

type fakeCustomers struct {
	customer customers.Customer
	address customers.Address
}

func (f fakeCustomers) GetCustomer(_ context.Context, tenantID, id string) (customers.Customer, bool, error) {
	if f.customer.TenantID == tenantID && f.customer.ID == id {
		return f.customer, true, nil
	}
	return customers.Customer{}, false, nil
}

func (f fakeCustomers) GetAddress(_ context.Context, tenantID, id string) (customers.Address, bool, error) {
	if f.address.TenantID == tenantID && f.address.ID == id {
		return f.address, true, nil
	}
	return customers.Address{}, false, nil
}

func TestServiceSnapshotsCustomerAndAddress(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{
		"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000},
	}}
	reader := fakeCustomers{
		customer: customers.Customer{TenantID: "tenant", ID: "customer-1", Name: "Cliente", Phone: "9999", Email: "c@example.com"},
		address: customers.Address{TenantID: "tenant", ID: "address-1", CustomerID: "customer-1", Street: "Rua A", Number: "10", Neighborhood: "Centro", City: "Santarém", State: "PA", PostalCode: "68000-000"},
	}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" }, reader)
	order, _, err := service.Create(context.Background(), "tenant", "idem-1", CreateInput{
		Source: "PDV", Items: []CreateLineInput{{ItemID: "item-1", Quantity: 1}},
		CustomerID: "customer-1", AddressID: "address-1", Fulfillment: FulfillmentDelivery,
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.Customer == nil || order.Customer.Name != "Cliente" || order.Address == nil || order.Address.Street != "Rua A" {
		t.Fatalf("missing snapshots: %+v", order)
	}
	if order.Fulfillment != FulfillmentDelivery {
		t.Fatalf("fulfillment = %s", order.Fulfillment)
	}
}

func TestServiceRejectsAddressFromAnotherCustomer(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000}}}
	reader := fakeCustomers{
		customer: customers.Customer{TenantID: "tenant", ID: "customer-1", Name: "Cliente", Phone: "9999"},
		address: customers.Address{TenantID: "tenant", ID: "address-1", CustomerID: "customer-2", Street: "Rua A", Number: "10", Neighborhood: "Centro", City: "Santarém", State: "PA", PostalCode: "68000-000"},
	}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" }, reader)
	_, _, err := service.Create(context.Background(), "tenant", "idem-1", CreateInput{Source: "PDV", Items: []CreateLineInput{{ItemID: "item-1", Quantity: 1}}, CustomerID: "customer-1", AddressID: "address-1"})
	if err != ErrAddressCustomerMismatch {
		t.Fatalf("err = %v, want address/customer mismatch", err)
	}
}

func TestServiceRequiresAddressForDelivery(t *testing.T) {
	catalog := &fakeCatalog{products: map[string]CatalogProduct{"item-1": {ID: "item-1", Name: "Pizza", PriceCents: 3000}}}
	reader := fakeCustomers{customer: customers.Customer{TenantID: "tenant", ID: "customer-1", Name: "Cliente", Phone: "9999"}}
	service := NewService(NewMemoryStore(), catalog, func() string { return "order-1" }, reader)
	_, _, err := service.Create(context.Background(), "tenant", "idem-1", CreateInput{Source: "PDV", Items: []CreateLineInput{{ItemID: "item-1", Quantity: 1}}, CustomerID: "customer-1", Fulfillment: FulfillmentDelivery})
	if err != ErrDeliveryAddressRequired {
		t.Fatalf("err = %v, want delivery address required", err)
	}
}
