package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"

    "github.com/TanyaKremnova/url-shortener/internal/auth"
    "github.com/TanyaKremnova/url-shortener/internal/models"
    "github.com/TanyaKremnova/url-shortener/internal/utils"
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
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Hash password — never store plain text
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not hash password")
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
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            utils.ErrorResponse(c, http.StatusConflict, "email already registered")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not create user")
        return
    }

    // Generate token right away so user is logged in after register
    token, err := auth.GenerateToken(userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not generate token")
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, models.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Find user by email
    var user models.User
    err := h.DB.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
        return
    }

    // Compare password with hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "could not generate token")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, models.AuthResponse{Token: token})
}