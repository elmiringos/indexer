package blockchain

import (
	"encoding/base64"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expected string
	}{
		{
			name:     "valid credentials",
			username: "user",
			password: "pass",
			expected: base64.StdEncoding.EncodeToString([]byte("user:pass")),
		},
		{
			name:     "empty credentials",
			username: "",
			password: "",
			expected: base64.StdEncoding.EncodeToString([]byte(":")),
		},
		{
			name:     "special characters",
			username: "test@user",
			password: "p@ssw0rd!",
			expected: base64.StdEncoding.EncodeToString([]byte("test@user:p@ssw0rd!")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := basicAuth(tt.username, tt.password)
			if result != tt.expected {
				t.Errorf("basicAuth(%q, %q) = %q; want %q", tt.username, tt.password, result, tt.expected)
			}
		})
	}
}
