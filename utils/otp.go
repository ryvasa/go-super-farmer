package utils

import (
	"crypto/rand"
	"fmt"
)

type OTP interface {
	GenerateOTP(length int) (string, error)
}
type OTPGenerator struct {
}

func NewOTPGenerator() OTP {
	return &OTPGenerator{}
}
func (g *OTPGenerator) GenerateOTP(length int) (string, error) {
	// Default length 6 if not specified
	if length == 0 {
		length = 6
	}

	// Generate random bytes
	numbers := make([]byte, length)
	if _, err := rand.Read(numbers); err != nil {
		return "", err
	}

	// Convert to numeric string
	var otp string
	for i := 0; i < length; i++ {
		otp += fmt.Sprintf("%d", numbers[i]%10)
	}

	return otp, nil
}
