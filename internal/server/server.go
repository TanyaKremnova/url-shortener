package server

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *gin.Engine {
    r := gin.Default()

    // Health check — always useful, mentor will like this
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
        })
    })

    // Route groups — stubs for now, handlers come in later tickets
    auth := r.Group("/auth")
    {
        auth.POST("/register", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
        auth.POST("/login", func(c *gin.Context) {
            c.JSON(501, gin.H{"message": "not implemented"})
        })
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

    // Redirect — must be last, catches /:code
    r.GET("/:code", func(c *gin.Context) {
        c.JSON(501, gin.H{"message": "not implemented"})
    })

    return r
}