package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "github.com/TanyaKremnova/url-shortener/internal/utils"
)

// Key used to store user_id in Gin context
const UserIDKey = "userID"

func Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header required")
            c.Abort() // Stop the request — don't call next handler
            return
        }

        // Header format must be: "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            utils.ErrorResponse(c, http.StatusUnauthorized, "authorization header format must be: Bearer <token>")
            c.Abort()
            return
        }

        // Validate the token
        claims, err := ValidateToken(parts[1])
        if err != nil {
            utils.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired token")
            c.Abort()
            return
        }

        // Store user_id in context so handlers can use it
        c.Set(UserIDKey, claims.UserID)

        c.Next()
    }
}