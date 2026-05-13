package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/models"
)

type AuthHandler struct {
    DB *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
    return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest

    // Validate request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password — never store plain text
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    // Insert user
    var userID string
    query := `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id
    `
    err = h.DB.QueryRowx(query, req.Email, string(hash)).Scan(&userID)
    if err != nil {
        // Postgres unique violation code = 23505
        if err.Error() != "" {
            c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
        return
    }

    // Generate token right away so user is logged in after register
    token, err := auth.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    c.JSON(http.StatusCreated, models.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find user by email
    var user models.User
    err := h.DB.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
    if err != nil {
        // Don't say "email not found" — always give the same message
        // so attackers can't enumerate valid emails
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Compare password with hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}