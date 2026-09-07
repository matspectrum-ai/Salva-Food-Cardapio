package identity

import (
	"context"
	"testing"
)

type fakeHasher struct {
	hashed string
	calls  int
}

func (f *fakeHasher) Hash(raw string) (string, error) {
	f.calls++
	return "test-hash:" + raw, nil
}
func TestServiceCreatesCollaborator(t *testing.T) {
	store := NewMemoryStore()
	hasher := &fakeHasher{}
	service := NewService(store, hasher, func() string { return "user-1" })
	input := CreateCollaboratorInput{
		Name: "Maria Silva", CPF: "529.982.247-25", Email: "MARIA@example.com",
		Phone: "(93) 99999-9999", Title: "Gerente",
		Permissions: []Permission{PermissionPOS, PermissionReports},
	}
	input.Password = "A1#aaaaa"
	created, err := service.CreateCollaborator(context.Background(), "tenant-a", input)
	if err != nil {
		t.Fatal(err)
	}
	if created.User.Email != "maria@example.com" || created.User.PasswordHash == "" || hasher.calls != 1 {
		t.Fatalf("created=%+v calls=%d", created, hasher.calls)
	}
}
func TestServiceRejectsWeakPasswordBeforeHashing(t *testing.T) {
	store := NewMemoryStore()
	hasher := &fakeHasher{}
	service := NewService(store, hasher, func() string { return "user-1" })
	input := CreateCollaboratorInput{
		Name: "Maria Silva", CPF: "529.982.247-25", Email: "maria@example.com",
		Phone: "(93) 99999-9999", Title: "Gerente",
		Permissions: []Permission{PermissionPOS},
	}
	input.Password = "weak"
	if _, err := service.CreateCollaborator(context.Background(), "tenant-a", input); err != ErrWeakPassword {
		t.Fatalf("err=%v want=%v", err, ErrWeakPassword)
	}
	if hasher.calls != 0 {
		t.Fatalf("hasher called %d times", hasher.calls)
	}
}
