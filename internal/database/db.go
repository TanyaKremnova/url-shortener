package database

import (
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // postgres driver, blank import registers it
)

func Connect(databaseURL string) *sqlx.DB {
    db, err := sqlx.Connect("postgres", databaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)

    log.Println("✅ Connected to database")
    return db
}