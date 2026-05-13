package handlers

import (
    "fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/models"
    "github.com/TanyaKremnova/url-shortener/internal/service"
)

type URLHandler struct {
    DB *sqlx.DB
}

func NewURLHandler(db *sqlx.DB) *URLHandler {
    return &URLHandler{DB: db}
}

func (h *URLHandler) CreateURL(c *gin.Context) {
    // Get user_id from context — middleware already validated the token
    userID, exists := c.Get(auth.UserIDKey)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    // Validate request body
    var req models.CreateURLRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Try to generate a unique short code
    // Retry up to 5 times in case of collision (extremely rare at small scale)
    var shortCode string
    var insertErr error

    for attempts := 0; attempts < 5; attempts++ {
        code, err := service.GenerateShortCode()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate short code"})
            return
        }

        query := `
            INSERT INTO urls (user_id, original_url, short_code)
            VALUES ($1, $2, $3)
            RETURNING short_code
        `
        insertErr = h.DB.QueryRowx(query, userID, req.OriginalURL, code).Scan(&shortCode)
        if insertErr == nil {
            break
        }

        // Check if it's a unique constraint violation (duplicate short_code)
        if pqErr, ok := insertErr.(*pq.Error); ok && pqErr.Code == "23505" {
            continue // try again with a new code
        }

        // Any other DB error
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save url"})
        return
    }

    if insertErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate unique code"})
        return
    }

    baseURL := os.Getenv("APP_BASE_URL")
    c.JSON(http.StatusCreated, models.CreateURLResponse{
        ShortCode:   shortCode,
        ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortCode),
        OriginalURL: req.OriginalURL,
    })
}