package identity

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordManager struct {
	Cost int
}

func NewBcryptPasswordManager(cost int) BcryptPasswordManager {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptPasswordManager{Cost: cost}
}

func (m BcryptPasswordManager) Hash(raw string) (string, error) {
	if err := ValidatePassword(raw); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), m.Cost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}
func (m BcryptPasswordManager) Compare(hash, raw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)); err != nil {
		return err
	}
	return nil
}

var _ PasswordHasher = BcryptPasswordManager{}
var _ PasswordVerifier = BcryptPasswordManager{}
