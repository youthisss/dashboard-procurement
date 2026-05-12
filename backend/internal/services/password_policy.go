package services

import (
	"errors"
	"strings"
)

var weakPasswords = map[string]bool{
	"admin123":                      true,
	"change_me":                     true,
	"password":                      true,
	"password123":                   true,
	"rygell_super_secret_change_me": true,
}

func ValidatePassword(password string) error {
	trimmed := strings.TrimSpace(password)
	if len(trimmed) < 12 {
		return errors.New("password must be at least 12 characters")
	}
	if weakPasswords[strings.ToLower(trimmed)] {
		return errors.New("password must not use a known placeholder")
	}

	hasLetter := false
	hasNumber := false
	for _, r := range trimmed {
		if r >= '0' && r <= '9' {
			hasNumber = true
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
	}
	if !hasLetter || !hasNumber {
		return errors.New("password must include letters and numbers")
	}
	return nil
}
