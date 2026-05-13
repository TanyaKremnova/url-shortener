package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

type RedirectHandler struct {
    DB *sqlx.DB
}

func NewRedirectHandler(db *sqlx.DB) *RedirectHandler {
    return &RedirectHandler{DB: db}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
    code := c.Param("code")

    // Look up original URL and increment click count atomically in one query
    // This is safer than SELECT then UPDATE separately (avoids race conditions)
    var originalURL string
    query := `
        UPDATE urls
        SET click_count = click_count + 1
        WHERE short_code = $1
        RETURNING original_url
    `
    err := h.DB.QueryRowx(query, code).Scan(&originalURL)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "short url not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
        return
    }

    // 302 = temporary redirect
    // NOT 301 — browsers cache 301 permanently, which would break our click counter
    c.Redirect(http.StatusFound, originalURL)
}