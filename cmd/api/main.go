package main

import (
    "log"

    "github.com/TanyaKremnova/url-shortener/internal/config"
    "github.com/TanyaKremnova/url-shortener/internal/database"
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

    log.Println("🚀 App started")
}