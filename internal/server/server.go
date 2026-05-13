package server

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/handlers"
    "github.com/TanyaKremnova/url-shortener/internal/middleware"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(middleware.Recovery())

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    authHandler := handlers.NewAuthHandler(db)
    urlHandler := handlers.NewURLHandler(db)
    redirectHandler := handlers.NewRedirectHandler(db)
    statsHandler := handlers.NewStatsHandler(db)

    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authHandler.Register)
        authGroup.POST("/login", authHandler.Login)
    }

    urls := r.Group("/urls", auth.Middleware())
    {
        urls.POST("/", urlHandler.CreateURL)
    }

    admin := r.Group("/admin", auth.Middleware())
    {
        admin.GET("/urls/stats", statsHandler.GetStats)
    }

    r.GET("/:code", redirectHandler.Redirect)

    return r
}