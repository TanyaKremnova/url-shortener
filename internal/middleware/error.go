package middleware

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the full error server-side
                log.Printf("PANIC: %v", err)

                // Return clean message to client — never expose internals
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "internal server error",
                    "code":  500,
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}