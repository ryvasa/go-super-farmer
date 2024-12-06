package utils

import (
	"golang.org/x/crypto/bcrypt"
)

var MockHashPassword func(password string) (string, error)

// HashPassword menggunakan MockHashPassword jika sudah di-set, atau menggunakan implementasi asli
func HashPassword(password string) (string, error) {
	if MockHashPassword != nil {
		return MockHashPassword(password) // menggunakan mock jika ada
	}

	// Fungsi asli jika mock tidak diset
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
