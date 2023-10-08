package handler

import (
	"testing"
)

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password123!", true},
		{"passw0rd", false},
		{"PASSWORD", false},
		{"password!", false},
		{"Password", false},
		{"Password123", false},
		{"Password!@#$", false},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			result := isValidPassword(test.password)
			if result != test.expected {
				t.Errorf("Password: %s, ожидалось: %v, получено: %v", test.password, test.expected, result)
			}
		})
	}
}
