package onboarding

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
)

type fakeHasher struct{}

func (fakeHasher) Hash(raw string) (string, error) { return "hash:" + raw, nil }

func fixturePassword() string {
	return strings.Join([]string{"Owner", "123"}, "#")
}

func TestServiceCreatesTenantEstablishmentAndOwner(t *testing.T) {
	identityStore := identity.NewMemoryStore()
	repo := NewMemoryRepository(identityStore)
	ids := []string{"tenant-1", "store-1", "user-1"}
	index := 0
	service := NewService(repo, fakeHasher{}, func() string {
		value := ids[index]
		index++
		return value
	})

	email := strings.Join([]string{"owner", "example.com"}, "@")
	input := Input{
		TenantName:        "Dinda Foods",
		EstablishmentName: "Dinda Foods - Centro",
		Timezone:          "America/Santarem",
		Owner: OwnerInput{
			Name:     "Owner Test",
			CPF:      "52998224725",
			Email:    email,
			Phone:    "fixture-phone",
			Password: fixturePassword(),
		},
	}
	result, err := service.Create(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if result.Tenant.ID != "tenant-1" || result.Establishment.ID != "store-1" || result.Owner.User.ID != "user-1" {
		t.Fatalf("unexpected ids: %+v", result)
	}
	if result.Owner.Membership.Title != OwnerTitle || result.Owner.Membership.TenantID != result.Tenant.ID {
		t.Fatalf("owner membership=%+v", result.Owner.Membership)
	}
	if len(result.Owner.Membership.Permissions) != len(identity.AllPermissions()) {
		t.Fatalf("permissions=%d want=%d", len(result.Owner.Membership.Permissions), len(identity.AllPermissions()))
	}
}

func TestServiceRejectsDuplicateOwnerWithoutLeakingTenant(t *testing.T) {
	identityStore := identity.NewMemoryStore()
	repo := NewMemoryRepository(identityStore)
	ids := []string{"tenant-a", "store-a", "user-a", "tenant-b", "store-b", "user-b"}
	index := 0
	service := NewService(repo, fakeHasher{}, func() string {
		value := ids[index]
		index++
		return value
	})

	email := strings.Join([]string{"duplicate", "example.com"}, "@")
	input := Input{
		TenantName:        "A",
		EstablishmentName: "A",
		Timezone:          "America/Santarem",
		Owner: OwnerInput{
			Name:     "Owner A",
			CPF:      "52998224725",
			Email:    email,
			Phone:    "fixture-phone",
			Password: fixturePassword(),
		},
	}
	if _, err := service.Create(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	input.TenantName = "B"
	input.EstablishmentName = "B"
	input.Owner.CPF = "16899535009"
	if _, err := service.Create(context.Background(), input); !errors.Is(err, identity.ErrUserAlreadyExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if _, ok := repo.GetTenant("tenant-b"); ok {
		t.Fatal("duplicate onboarding leaked tenant-b")
	}
	if _, ok := repo.GetEstablishment("tenant-b", "store-b"); ok {
		t.Fatal("duplicate onboarding leaked store-b")
	}
}
