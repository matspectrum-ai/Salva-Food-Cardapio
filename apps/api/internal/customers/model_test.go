package customers

import "testing"

func TestNewCustomerValidatesRequiredFields(t *testing.T) {
	if _, err := NewCustomer("", "c1", "Ana", "91999999999", ""); err != ErrTenantRequired {
		t.Fatal(err)
	}
	if _, err := NewCustomer("t1", "", "Ana", "91999999999", ""); err != ErrIDRequired {
		t.Fatal(err)
	}
	if _, err := NewCustomer("t1", "c1", "", "91999999999", ""); err != ErrNameRequired {
		t.Fatal(err)
	}
	if _, err := NewCustomer("t1", "c1", "Ana", "", ""); err != ErrPhoneRequired {
		t.Fatal(err)
	}
}

func TestNewCustomerTrimsFields(t *testing.T) {
	c, err := NewCustomer("t1", "c1", " Ana ", " 91999999999 ", " a@example.com ")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "Ana" || c.Phone != "91999999999" || c.Email != "a@example.com" {
		t.Fatalf("unexpected customer: %+v", c)
	}
}

func TestNewAddressRequiresCustomer(t *testing.T) {
	if _, err := NewAddress("t1", "a1", "", "Rua A", "10", "Centro", "Santarém", "PA", "68000-000"); err != ErrCustomerRequired {
		t.Fatal(err)
	}
}
