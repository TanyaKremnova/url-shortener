package server

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "github.com/TanyaKremnova/url-shortener/internal/handlers"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    authHandler := handlers.NewAuthHandler(db)

    auth := r.Group("/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
    }

    urls := r.Group("/urls")
    {
        urls.POST("/", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
    }

    admin := r.Group("/admin")
    {
        admin.GET("/urls/stats", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
    }

    r.GET("/:code", func(c *gin.Context) {
        c.JSON(501, gin.H{"message": "not implemented"})
    })

    return r
}