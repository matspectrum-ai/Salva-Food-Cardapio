package identity

import (
	"context"
	"testing"
)

func validCollaborator(t *testing.T, tenantID, userID, email string) Collaborator {
	t.Helper()
	user, err := NewUser(userID, "Maria Silva", "529.982.247-25", email, "(93) 99999-9999", "hash")
	if err != nil {
		t.Fatal(err)
	}
	membership, err := NewMembership(tenantID, userID, "Gerente", []Permission{PermissionPOS, PermissionReports})
	if err != nil {
		t.Fatal(err)
	}
	return Collaborator{User: user, Membership: membership}
}

func TestMemoryStoreCreatesAndListsCollaborators(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	collaborator := validCollaborator(t, "tenant-a", "user-1", "maria@example.com")
	if err := store.CreateCollaborator(ctx, collaborator); err != nil {
		t.Fatal(err)
	}
	items, err := store.ListCollaborators(ctx, "tenant-a", "gerente")
	if err != nil || len(items) != 1 || items[0].User.Email != "maria@example.com" {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}
func TestMemoryStoreRejectsDuplicateIdentity(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	first := validCollaborator(t, "tenant-a", "user-1", "maria@example.com")
	second := validCollaborator(t, "tenant-b", "user-2", "maria@example.com")
	if err := store.CreateCollaborator(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCollaborator(ctx, second); err != ErrUserAlreadyExists {
		t.Fatalf("err=%v want=%v", err, ErrUserAlreadyExists)
	}
}

func TestMemoryStoreIsolatesMembershipsByTenant(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	collaborator := validCollaborator(t, "tenant-a", "user-1", "maria@example.com")
	if err := store.CreateCollaborator(ctx, collaborator); err != nil {
		t.Fatal(err)
	}
	items, err := store.ListCollaborators(ctx, "tenant-b", "")
	if err != nil || len(items) != 0 {
		t.Fatalf("tenant leak items=%+v err=%v", items, err)
	}
}
func TestMemoryStoreChangesMembershipStatus(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	collaborator := validCollaborator(t, "tenant-a", "user-1", "maria@example.com")
	if err := store.CreateCollaborator(ctx, collaborator); err != nil {
		t.Fatal(err)
	}
	if err := store.SetMembershipStatus(ctx, "tenant-a", "user-1", MembershipInactive); err != nil {
		t.Fatal(err)
	}
	membership, ok, err := store.GetMembership(ctx, "tenant-a", "user-1")
	if err != nil || !ok || membership.Status != MembershipInactive {
		t.Fatalf("membership=%+v ok=%v err=%v", membership, ok, err)
	}
	if err := store.SetMembershipStatus(ctx, "tenant-b", "user-1", MembershipInactive); err != ErrCollaboratorNotFound {
		t.Fatalf("err=%v want=%v", err, ErrCollaboratorNotFound)
	}
}
