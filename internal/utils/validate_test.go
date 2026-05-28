package utils

import "testing"

func TestIsValidURL(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"https://www.google.com", true},
        {"http://example.com", true},
        {"http://localhost:8080", true},
        {"not-a-url", false},
        {"ftp://example.com", false},
        {"", false},
        {"https://", false},
        {"javascript:alert(1)", false},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := IsValidURL(tt.input)
            if result != tt.expected {
                t.Errorf("IsValidURL(%q) = %v, want %v", tt.input, result, tt.expected)
            }
        })
    }
}