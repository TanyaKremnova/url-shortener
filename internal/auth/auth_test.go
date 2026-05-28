package auth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/TanyaKremnova/url-shortener/internal/auth"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func setSecret(t *testing.T, secret string) {
	t.Helper()
	os.Setenv("JWT_SECRET", secret)
	t.Cleanup(func() { os.Unsetenv("JWT_SECRET") })
}

// makeExpiredToken creates a token that is already expired — for testing
func makeExpiredToken(t *testing.T, userID string) string {
	t.Helper()
	claims := auth.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // past
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("could not create expired token: %v", err)
	}
	return signed
}

// ── JWT tests ─────────────────────────────────────────────────────────────────

func TestGenerateToken_ReturnsToken(t *testing.T) {
	setSecret(t, "test-secret")

	token, err := auth.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestGenerateToken_DifferentUsersGetDifferentTokens(t *testing.T) {
	setSecret(t, "test-secret")

	token1, _ := auth.GenerateToken("user-1")
	token2, _ := auth.GenerateToken("user-2")

	if token1 == token2 {
		t.Error("different users should get different tokens")
	}
}

func TestValidateToken_ValidToken(t *testing.T) {
	setSecret(t, "test-secret")

	token, err := auth.GenerateToken("user-abc")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.UserID != "user-abc" {
		t.Errorf("expected user_id 'user-abc', got '%s'", claims.UserID)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Generate with one secret
	setSecret(t, "secret-A")
	token, _ := auth.GenerateToken("user-123")

	// Validate with a different secret
	os.Setenv("JWT_SECRET", "secret-B")

	_, err := auth.ValidateToken(token)
	if err == nil {
		t.Error("expected error when validating with wrong secret, got nil")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	setSecret(t, "test-secret")

	expired := makeExpiredToken(t, "user-123")

	_, err := auth.ValidateToken(expired)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateToken_TamperedToken(t *testing.T) {
	setSecret(t, "test-secret")

	token, _ := auth.GenerateToken("user-123")

	// Flip the last character — breaks the signature
	tampered := token[:len(token)-1] + "X"

	_, err := auth.ValidateToken(tampered)
	if err == nil {
		t.Error("expected error for tampered token, got nil")
	}
}

func TestValidateToken_EmptyString(t *testing.T) {
	setSecret(t, "test-secret")

	_, err := auth.ValidateToken("")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

func TestValidateToken_RandomString(t *testing.T) {
	setSecret(t, "test-secret")

	_, err := auth.ValidateToken("this.is.not.a.jwt")
	if err == nil {
		t.Error("expected error for random string, got nil")
	}
}

// ── Middleware tests ──────────────────────────────────────────────────────────

// runMiddleware fires the auth middleware against a fake request
// and returns the response recorder so we can inspect status + body
func runMiddleware(t *testing.T, authHeader string) (*httptest.ResponseRecorder, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, engine := gin.CreateTestContext(w)

	// Register the middleware + a dummy handler that marks "reached"
	reached := false
	engine.GET("/test", auth.Middleware(), func(c *gin.Context) {
		reached = true
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	c.Request = req

	engine.ServeHTTP(w, req)
	return w, reached
}

func TestMiddleware_NoHeader(t *testing.T) {
	setSecret(t, "test-secret")

	w, reached := runMiddleware(t, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if reached {
		t.Error("handler should not have been reached")
	}
}

func TestMiddleware_WrongFormat_NoBearer(t *testing.T) {
	setSecret(t, "test-secret")

	w, reached := runMiddleware(t, "Basic sometoken")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if reached {
		t.Error("handler should not have been reached")
	}
}

func TestMiddleware_WrongFormat_BearerOnly(t *testing.T) {
	setSecret(t, "test-secret")

	// "Bearer" with no token after it
	w, reached := runMiddleware(t, "Bearer")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if reached {
		t.Error("handler should not have been reached")
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	setSecret(t, "test-secret")

	w, reached := runMiddleware(t, "Bearer thisisnotavalidtoken")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if reached {
		t.Error("handler should not have been reached")
	}
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	setSecret(t, "test-secret")

	expired := makeExpiredToken(t, "user-123")
	w, reached := runMiddleware(t, "Bearer "+expired)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if reached {
		t.Error("handler should not have been reached")
	}
}

func TestMiddleware_ValidToken_PassesThrough(t *testing.T) {
	setSecret(t, "test-secret")

	token, _ := auth.GenerateToken("user-xyz")
	w, reached := runMiddleware(t, "Bearer "+token)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
	if !reached {
		t.Error("handler should have been reached with valid token")
	}
}

func TestMiddleware_ValidToken_SetsUserIDInContext(t *testing.T) {
	setSecret(t, "test-secret")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)

	var capturedUserID string

	engine.GET("/test", auth.Middleware(), func(c *gin.Context) {
		// Read user_id that middleware stored in context
		id, exists := c.Get(auth.UserIDKey)
		if !exists {
			t.Error("userID not found in context")
			return
		}
		capturedUserID = id.(string)
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token, _ := auth.GenerateToken("user-context-check")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	engine.ServeHTTP(w, req)

	if capturedUserID != "user-context-check" {
		t.Errorf("expected userID 'user-context-check' in context, got '%s'", capturedUserID)
	}
}