package service

import (
    "testing"
)

func TestGenerateShortCode_Length(t *testing.T) {
    code, err := GenerateShortCode()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(code) != 6 {
        t.Errorf("expected length 6, got %d", len(code))
    }
}

func TestGenerateShortCode_Uniqueness(t *testing.T) {
    codes := make(map[string]bool)
    for i := 0; i < 1000; i++ {
        code, err := GenerateShortCode()
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if codes[code] {
            t.Errorf("duplicate code generated: %s", code)
        }
        codes[code] = true
    }
}

func TestGenerateShortCode_ValidChars(t *testing.T) {
    const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    for i := 0; i < 100; i++ {
        code, _ := GenerateShortCode()
        for _, ch := range code {
            found := false
            for _, v := range validChars {
                if ch == v {
                    found = true
                    break
                }
            }
            if !found {
                t.Errorf("invalid character in code: %c", ch)
            }
        }
    }
}