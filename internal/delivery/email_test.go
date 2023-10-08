package handler

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"user@example.com", true},
		{"user123@exampleco.uk", true},
		{"invalid-email", false},
		{"user@.com", false},
		{"user@.example.com", false},
		{"user@example", false},
		{"user.example.com", false},
		{"user@.example", false},
		{"user.example", false},
		{"user@example@com", false},
		{"user@example@.com", false},
		{"user@example.com@.com", false},
		{"user@.example@.com", false},
		{"user@example.com@.example.com", false},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			result := isValidEmail(test.email)
			if result != test.expected {
				t.Errorf("Email: %s, ожидалось: %v, получено: %v", test.email, test.expected, result)
			}
		})
	}
}
