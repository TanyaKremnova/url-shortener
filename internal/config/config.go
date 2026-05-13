package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    AppPort     string
    JWTSecret   string
    AppBaseURL  string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading from environment")
    }

    return &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        AppPort:     os.Getenv("APP_PORT"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
        AppBaseURL:  os.Getenv("APP_BASE_URL"),
    }
}