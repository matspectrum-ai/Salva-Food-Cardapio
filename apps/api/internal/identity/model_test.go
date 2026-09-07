package identity

import "testing"

func TestPasswordPolicy(t *testing.T) {
	valid := []string{
		"Strong#1",
		"Outra@Senha2",
		"Senha$MuitoForte9",
	}
	for _, password := range valid {
		if err := ValidatePassword(password); err != nil {
			t.Fatalf("valid password %q rejected: %v", password, err)
		}
	}

	invalid := []string{
		"Short#1",
		"lowercase#1",
		"UPPERCASE#1",
		"NoNumber#",
		"NoSymbol1A",
	}
	for _, password := range invalid {
		if err := ValidatePassword(password); err == nil {
			t.Fatalf("invalid password %q accepted", password)
		}
	}
}
func TestPermissionSetIsValidatedAndCanonicalized(t *testing.T) {
	permissions, err := NormalizePermissionSet([]Permission{
		PermissionReports,
		PermissionPOS,
		PermissionReports,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(permissions) != 2 {
		t.Fatalf("permissions len = %d, want 2", len(permissions))
	}
	if permissions[0] != PermissionPOS || permissions[1] != PermissionReports {
		t.Fatalf("unexpected canonical order: %+v", permissions)
	}

	if _, err := NormalizePermissionSet([]Permission{"unknown.permission"}); err != ErrInvalidPermission {
		t.Fatalf("err = %v, want %v", err, ErrInvalidPermission)
	}
	if _, err := NormalizePermissionSet(nil); err != ErrPermissionsRequired {
		t.Fatalf("err = %v, want %v", err, ErrPermissionsRequired)
	}
}

func TestNewUserNormalizesIdentity(t *testing.T) {
	user, err := NewUser(
		"11111111-1111-4111-8111-111111111111",
		"  Maria Silva  ",
		"529.982.247-25",
		"  MARIA@EXAMPLE.COM ",
		"(93) 99999-9999",
		"hash",
	)
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != "Maria Silva" || user.CPF != "52998224725" || user.Email != "maria@example.com" {
		t.Fatalf("unexpected normalized user: %+v", user)
	}
	if user.PasswordHash != "hash" {
		t.Fatal("password hash was not preserved")
	}
}
func TestNewUserRejectsInvalidCPFAndEmail(t *testing.T) {
	if _, err := NewUser("id", "Maria", "111.111.111-11", "maria@example.com", "phone", "hash"); err != ErrInvalidCPF {
		t.Fatalf("cpf err = %v, want %v", err, ErrInvalidCPF)
	}
	if _, err := NewUser("id", "Maria", "529.982.247-25", "not-an-email", "phone", "hash"); err != ErrInvalidEmail {
		t.Fatalf("email err = %v, want %v", err, ErrInvalidEmail)
	}
}

func TestMembershipSeparatesTitleFromAuthorization(t *testing.T) {
	membership, err := NewMembership(
		"22222222-2222-4222-8222-222222222222",
		"11111111-1111-4111-8111-111111111111",
		"  Gerente de turno  ",
		[]Permission{PermissionPOS, PermissionReports},
	)
	if err != nil {
		t.Fatal(err)
	}
	if membership.Title != "Gerente de turno" || membership.Status != MembershipActive {
		t.Fatalf("unexpected membership: %+v", membership)
	}
	if !membership.HasPermission(PermissionPOS) || membership.HasPermission(PermissionFiscalIssue) {
		t.Fatalf("permission set incorrect: %+v", membership.Permissions)
	}
	membership.SetStatus(MembershipInactive)
	if membership.Status != MembershipInactive {
		t.Fatalf("status = %s", membership.Status)
	}
}
