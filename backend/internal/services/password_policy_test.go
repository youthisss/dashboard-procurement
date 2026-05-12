package services

import "testing"

func TestValidatePasswordRejectsWeakPasswords(t *testing.T) {
	tests := []string{
		"short1",
		"admin123",
		"change_me",
		"averylongpassword",
		"123456789012",
	}

	for _, password := range tests {
		t.Run(password, func(t *testing.T) {
			if err := ValidatePassword(password); err == nil {
				t.Fatalf("ValidatePassword(%q) returned nil, want error", password)
			}
		})
	}
}

func TestValidatePasswordAllowsStrongPassword(t *testing.T) {
	if err := ValidatePassword("Procurement2026Strong"); err != nil {
		t.Fatalf("ValidatePassword() error = %v", err)
	}
}
