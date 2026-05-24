package handlers

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/cache"
    "github.com/TanyaKremnova/url-shortener/internal/utils"
)

type RedirectHandler struct {
    DB *sqlx.DB
    Cache *cache.Cache
}

func NewRedirectHandler(db *sqlx.DB, cache *cache.Cache) *RedirectHandler {
    return &RedirectHandler{
        DB:     db,
        Cache:  cache,}
}

func (h *RedirectHandler) incrementClickCount(code string) {
    query := `
        UPDATE urls
        SET click_count = click_count + 1
        WHERE short_code = $1
    `
    if _, err := h.DB.Exec(query, code); err != nil {
        log.Printf("failed to increment click count for %s: %v", code, err)
    }
    h.DB.Exec(query, code)
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
    code := c.Param("code")
    ctx := c.Request.Context()

    // 1. Check cache first
    if originalURL, err := h.Cache.Get(ctx, code); err == nil {
        // Cache hit — no DB query needed
        go h.incrementClickCount(code) //TODO explain concurrency usage
        c.Redirect(http.StatusFound, originalURL)
        return
    }

    // 2. Cache miss — query DB
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

    // 3. Store in cache for next time
    h.Cache.Set(ctx, code, originalURL)

    c.Redirect(http.StatusFound, originalURL)
}