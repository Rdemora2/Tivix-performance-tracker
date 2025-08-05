package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
)

func GenerateTemporaryPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*"
	const length = 12

	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("senha deve ter pelo menos 8 caracteres")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper {
		return errors.New("senha deve conter pelo menos uma letra maiúscula")
	}

	if !hasLower {
		return errors.New("senha deve conter pelo menos uma letra minúscula")
	}

	if !hasNumber {
		return errors.New("senha deve conter pelo menos um número")
	}

	return nil
}
