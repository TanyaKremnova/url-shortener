package main

import (
    "log"

    "github.com/TanyaKremnova/url-shortener/internal/config"
    "github.com/TanyaKremnova/url-shortener/internal/database"
    "github.com/TanyaKremnova/url-shortener/internal/server"
    "github.com/TanyaKremnova/url-shortener/internal/cache"
)

func main() {
    // Load config
    cfg := config.Load()

    if cfg.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    // Connect to DB
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()

    c := cache.New(cfg.CacheURL)

    r := server.NewRouter(db, c)

    log.Printf("🚀 Server running on port %s", cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}