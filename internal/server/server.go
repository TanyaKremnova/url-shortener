package server

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/handlers"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    authHandler := handlers.NewAuthHandler(db)

    // Public routes
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authHandler.Register)
        authGroup.POST("/login", authHandler.Login)
    }

    // Protected routes — middleware runs before every handler in these groups
    urls := r.Group("/urls", auth.Middleware())
    {
        urls.POST("/", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
    }

    admin := r.Group("/admin", auth.Middleware())
    {
        admin.GET("/urls/stats", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
    }

    // Public — redirect
    r.GET("/:code", func(c *gin.Context) {
        c.JSON(501, gin.H{"message": "not implemented"})
    })

    return r
}