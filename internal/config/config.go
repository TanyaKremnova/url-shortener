package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    AppPort     string
}

func Load() *Config {
    // In production this does nothing (env vars already set)
    // Locally it loads from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading from environment")
    }

    return &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        AppPort:     os.Getenv("APP_PORT"),
    }
}