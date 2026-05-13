package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/models"
    "github.com/TanyaKremnova/url-shortener/internal/utils"
)

type StatsHandler struct {
    DB *sqlx.DB
}

func NewStatsHandler(db *sqlx.DB) *StatsHandler {
    return &StatsHandler{DB: db}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
    // Read user_id from context — set by JWT middleware
    // NEVER take user_id from the request body or query params
    // Always trust the token, never the client
    userID, exists := c.Get(auth.UserIDKey)
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    var urls []models.URLStats

    query := `
        SELECT short_code, original_url, click_count, created_at
        FROM urls
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    err := h.DB.Select(&urls, query, userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not fetch stats")
        return
    }

    // Return empty array not null when user has no URLs
    if urls == nil {
        urls = []models.URLStats{}
    }

    utils.SuccessResponse(c, http.StatusOK, models.StatsResponse{
        URLs:  urls,
        Total: len(urls),
    })
}