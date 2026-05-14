package utils

import "github.com/gin-gonic/gin"

// Standard error shape: {"error": "message", "code": 404}
func ErrorResponse(c *gin.Context, status int, message string) {
    c.JSON(status, gin.H{
        "error": message,
        "code":  status,
    })
}

// Standard success shape: {"data": ..., "code": 200}
func SuccessResponse(c *gin.Context, status int, data interface{}) {
    c.JSON(status, gin.H{
        "data": data,
        "code": status,
    })
}