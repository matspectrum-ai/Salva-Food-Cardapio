package tenant

import (
	"errors"
	"testing"
)

func TestTenantAndEstablishmentValidation(t *testing.T) {
	tenant, err := New("tenant-1", "Restaurante Central")
	if err != nil {
		t.Fatal(err)
	}
	if tenant.Name != "Restaurante Central" {
		t.Fatalf("name=%q", tenant.Name)
	}

	establishment, err := NewEstablishment(
		"store-1", tenant.ID, "Unidade Centro", "America/Santarem",
	)
	if err != nil {
		t.Fatal(err)
	}
	if establishment.Timezone != "America/Santarem" {
		t.Fatalf("timezone=%q", establishment.Timezone)
	}

	if _, err := NewEstablishment(
		"store-2", tenant.ID, "Unidade", "Mars/Olympus",
	); !errors.Is(err, ErrInvalidTimezone) {
		t.Fatalf("invalid timezone err=%v", err)
	}
}

func TestTenantRequiresIdentityAndName(t *testing.T) {
	if _, err := New("", "Restaurante"); !errors.Is(err, ErrIDRequired) {
		t.Fatalf("id err=%v", err)
	}
	if _, err := New("tenant-1", "  "); !errors.Is(err, ErrNameRequired) {
		t.Fatalf("name err=%v", err)
	}
	if _, err := NewEstablishment(
		"store-1", "", "Unidade", "America/Sao_Paulo",
	); !errors.Is(err, ErrTenantRequired) {
		t.Fatalf("tenant err=%v", err)
	}
}
