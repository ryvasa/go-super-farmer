package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password, hash string) bool
}

type HasherImpl struct{}

func NewHasher() Hasher {
	return &HasherImpl{}
}

func (p *HasherImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (p *HasherImpl) ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
