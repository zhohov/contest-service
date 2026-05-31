package infrastructure

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// type DummyPasswordHasher struct{}

// func NewDummyPasswordHasher() *DummyPasswordHasher {
// 	return &DummyPasswordHasher{}
// }

// func (h *DummyPasswordHasher) Hash(password string) (string, error) {
// 	return password, nil
// }

// func (h *DummyPasswordHasher) Compare(password, hash string) bool {
// 	return password == hash
// }

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}

	return string(hash), nil
}

func (h *Hasher) Compare(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
