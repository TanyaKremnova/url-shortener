package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/utils"
)

type RedirectHandler struct {
    DB *sqlx.DB
}

func NewRedirectHandler(db *sqlx.DB) *RedirectHandler {
    return &RedirectHandler{DB: db}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
    code := c.Param("code")

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
            utils.ErrorResponse(c, http.StatusNotFound, "short url not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "something went wrong")
        return
    }

    c.Redirect(http.StatusFound, originalURL)
}