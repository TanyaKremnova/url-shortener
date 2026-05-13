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
    "github.com/TanyaKremnova/url-shortener/internal/utils"
)

type URLHandler struct {
    DB *sqlx.DB
}

func NewURLHandler(db *sqlx.DB) *URLHandler {
    return &URLHandler{DB: db}
}

func (h *URLHandler) CreateURL(c *gin.Context) {
    // Get user_id from context
    userID, exists := c.Get(auth.UserIDKey)
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    // Validate request body
    var req models.CreateURLRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Extra URL safety check
    if !utils.IsValidURL(req.OriginalURL) {
        utils.ErrorResponse(c, http.StatusBadRequest, "invalid url: must start with http:// or https://")
        return
    }

    // Try to generate a unique short code
    // Retry up to 5 times in case of collision (extremely rare at small scale)
    var shortCode string
    var insertErr error

    for attempts := 0; attempts < 5; attempts++ {
        code, err := service.GenerateShortCode()
        if err != nil {
            utils.ErrorResponse(c, http.StatusInternalServerError, "could not generate short code")
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
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not save url")
        return
    }

    if insertErr != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not generate unique code")
        return
    }

    baseURL := os.Getenv("APP_BASE_URL")
    utils.SuccessResponse(c, http.StatusCreated, models.CreateURLResponse{
        ShortCode:   shortCode,
        ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortCode),
        OriginalURL: req.OriginalURL,
    })
}