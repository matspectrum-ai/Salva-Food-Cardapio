package identity

import "testing"

func TestBcryptPasswordManager(t *testing.T) {
	manager := NewBcryptPasswordManager(4)
	hash, err := manager.Hash("Strong#123")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "Strong#123" || hash == "" {
		t.Fatalf("unexpected hash %q", hash)
	}
	if err := manager.Compare(hash, "Strong#123"); err != nil {
		t.Fatalf("compare valid password: %v", err)
	}
	if err := manager.Compare(hash, "Wrong#123"); err == nil {
		t.Fatal("wrong password unexpectedly matched")
	}
}

func TestBcryptRejectsWeakPassword(t *testing.T) {
	manager := NewBcryptPasswordManager(4)
	if _, err := manager.Hash("weak"); err != ErrWeakPassword {
		t.Fatalf("err=%v want=%v", err, ErrWeakPassword)
	}
}
