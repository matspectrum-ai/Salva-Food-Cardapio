package customers

import (
	"errors"
	"strings"
)

var (
	ErrTenantRequired   = errors.New("tenant id is required")
	ErrIDRequired       = errors.New("id is required")
	ErrNameRequired     = errors.New("name is required")
	ErrPhoneRequired    = errors.New("phone is required")
	ErrCustomerRequired = errors.New("customer id is required")
)

type Customer struct {
	TenantID string `json:"tenant_id"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email,omitempty"`
}

func NewCustomer(tenantID, id, name, phone, email string) (Customer, error) {
	if strings.TrimSpace(tenantID) == "" {
		return Customer{}, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return Customer{}, ErrIDRequired
	}
	if strings.TrimSpace(name) == "" {
		return Customer{}, ErrNameRequired
	}
	if strings.TrimSpace(phone) == "" {
		return Customer{}, ErrPhoneRequired
	}
	return Customer{TenantID: tenantID, ID: id, Name: strings.TrimSpace(name), Phone: strings.TrimSpace(phone), Email: strings.TrimSpace(email)}, nil
}

type Address struct {
	TenantID     string `json:"tenant_id"`
	ID           string `json:"id"`
	CustomerID   string `json:"customer_id"`
	Label        string `json:"label,omitempty"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Reference    string `json:"reference,omitempty"`
}

func NewAddress(tenantID, id, customerID, street, number, neighborhood, city, state, postalCode string) (Address, error) {
	if strings.TrimSpace(tenantID) == "" {
		return Address{}, ErrTenantRequired
	}
	if strings.TrimSpace(id) == "" {
		return Address{}, ErrIDRequired
	}
	if strings.TrimSpace(customerID) == "" {
		return Address{}, ErrCustomerRequired
	}
	return Address{TenantID: tenantID, ID: id, CustomerID: customerID, Street: strings.TrimSpace(street), Number: strings.TrimSpace(number), Neighborhood: strings.TrimSpace(neighborhood), City: strings.TrimSpace(city), State: strings.TrimSpace(state), PostalCode: strings.TrimSpace(postalCode)}, nil
}
