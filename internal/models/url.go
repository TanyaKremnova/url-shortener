package models

import "time"

type URL struct {
    ID          string    `db:"id"`
    UserID      string    `db:"user_id"`
    OriginalURL string    `db:"original_url"`
    ShortCode   string    `db:"short_code"`
    ClickCount  int       `db:"click_count"`
    CreatedAt   time.Time `db:"created_at"`
}

// Accept from the request body
type CreateURLRequest struct {
    OriginalURL string `json:"original_url" binding:"required,url"`
}

// Send back
type CreateURLResponse struct {
    ShortCode   string `json:"short_code"`
    ShortURL    string `json:"short_url"`
    OriginalURL string `json:"original_url"`
}