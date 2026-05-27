package models

import (
    "encoding/json"
    "testing"
)

func TestRegisterRequest_JSON(t *testing.T) {
    request := RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
    }

    data, err := json.Marshal(request)
    if err != nil {
        t.Fatalf("failed to marshal request: %v", err)
    }

    expected := `{"email":"test@example.com","password":"password123"}`

    if string(data) != expected {
        t.Errorf("expected %s, got %s", expected, string(data))
    }
}

func TestLoginRequest_JSON(t *testing.T) {
    request := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }

    data, err := json.Marshal(request)
    if err != nil {
        t.Fatalf("failed to marshal request: %v", err)
    }

    expected := `{"email":"test@example.com","password":"password123"}`

    if string(data) != expected {
        t.Errorf("expected %s, got %s", expected, string(data))
    }
}

func TestAuthResponse_JSON(t *testing.T) {
    response := AuthResponse{
        Token: "jwt-token",
    }

    data, err := json.Marshal(response)
    if err != nil {
        t.Fatalf("failed to marshal response: %v", err)
    }

    expected := `{"token":"jwt-token"}`

    if string(data) != expected {
        t.Errorf("expected %s, got %s", expected, string(data))
    }
}