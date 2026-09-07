package identity

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakePasswordManager struct{}

func (fakePasswordManager) Hash(raw string) (string, error) {
	return "test-hash:" + raw, nil
}
func (fakePasswordManager) Compare(hash, raw string) error {
	if hash != "test-hash:"+raw {
		return errors.New("mismatch")
	}
	return nil
}

func authFixture(t *testing.T) (*MemoryStore, *AuthService, time.Time) {
	t.Helper()
	store := NewMemoryStore()
	user, err := NewUser("user-1", "Maria Silva", "529.982.247-25", "maria@example.com", "phone", "test-hash:A1#aaaaa")
	if err != nil {
		t.Fatal(err)
	}
	membership, err := NewMembership("tenant-a", "user-1", "Gerente", []Permission{PermissionPOS, PermissionReports})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateCollaborator(context.Background(), Collaborator{User: user, Membership: membership}); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC)
	service := NewAuthService(store, fakePasswordManager{}, func() (string, error) { return "raw-session-token", nil }, func() string { return "session-1" }, func() time.Time { return now }, 12*time.Hour)
	return store, service, now
}
func TestAuthLoginAuthenticateAndLogout(t *testing.T) {
	_, service, now := authFixture(t)
	result, err := service.Login(context.Background(), LoginInput{
		Email: "MARIA@example.com", Password: "A1#aaaaa",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Token != "raw-session-token" || !result.ExpiresAt.Equal(now.Add(12*time.Hour)) {
		t.Fatalf("unexpected login result: %+v", result)
	}
	principal, err := service.Authenticate(context.Background(), result.Token)
	if err != nil {
		t.Fatal(err)
	}
	if principal.User.Email != "maria@example.com" || !principal.Membership.HasPermission(PermissionPOS) {
		t.Fatalf("unexpected principal: %+v", principal)
	}
	if err := service.Logout(context.Background(), result.Token); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Authenticate(context.Background(), result.Token); err != ErrSessionRevoked {
		t.Fatalf("err=%v want=%v", err, ErrSessionRevoked)
	}
}
func TestAuthRejectsInvalidCredentialsAndInactiveNewLogin(t *testing.T) {
	store, service, _ := authFixture(t)
	if _, err := service.Login(context.Background(), LoginInput{Email: "maria@example.com", Password: "wrong"}); err != ErrInvalidCredentials {
		t.Fatalf("invalid credentials err=%v", err)
	}
	if err := store.SetMembershipStatus(context.Background(), "tenant-a", "user-1", MembershipInactive); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Login(context.Background(), LoginInput{Email: "maria@example.com", Password: "A1#aaaaa"}); err != ErrNoActiveMembership {
		t.Fatalf("inactive login err=%v", err)
	}
}

func TestExistingSessionSurvivesMembershipDeactivation(t *testing.T) {
	store, service, _ := authFixture(t)
	result, err := service.Login(context.Background(), LoginInput{Email: "maria@example.com", Password: "A1#aaaaa"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetMembershipStatus(context.Background(), "tenant-a", "user-1", MembershipInactive); err != nil {
		t.Fatal(err)
	}
	principal, err := service.Authenticate(context.Background(), result.Token)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Membership.Status != MembershipInactive {
		t.Fatalf("status=%s want=%s", principal.Membership.Status, MembershipInactive)
	}
}
