package models

import (
    "encoding/json"
    "testing"
    "time"
)

func TestCreateURLRequest_JSON(t *testing.T) {
    request := CreateURLRequest{
        OriginalURL: "https://github.com",
    }

    data, err := json.Marshal(request)
    if err != nil {
        t.Fatalf("failed to marshal request: %v", err)
    }

    expected := `{"original_url":"https://github.com"}`

    if string(data) != expected {
        t.Errorf("expected %s, got %s", expected, string(data))
    }
}

func TestCreateURLResponse_JSON(t *testing.T) {
    response := CreateURLResponse{
        ShortCode:   "abc123",
        ShortURL:    "http://localhost/abc123",
        OriginalURL: "https://github.com",
    }

    data, err := json.Marshal(response)
    if err != nil {
        t.Fatalf("failed to marshal response: %v", err)
    }

    expected := `{"short_code":"abc123","short_url":"http://localhost/abc123","original_url":"https://github.com"}`

    if string(data) != expected {
        t.Errorf("expected %s, got %s", expected, string(data))
    }
}

func TestStatsResponse(t *testing.T) {
    stats := StatsResponse{
        URLs: []URLStats{
            {
                ShortCode:   "abc123",
                OriginalURL: "https://github.com",
                ClickCount:  10,
                CreatedAt:   time.Now(),
            },
        },
        Total: 1,
    }

    if stats.Total != 1 {
        t.Errorf("expected total 1")
    }

    if len(stats.URLs) != 1 {
        t.Errorf("expected 1 url")
    }

    if stats.URLs[0].ClickCount != 10 {
        t.Errorf("expected click count 10")
    }
}