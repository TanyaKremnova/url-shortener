package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecovery_NoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Recovery())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRecovery_WithPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Recovery())

	r.GET("/panic", func(c *gin.Context) {
		panic("something failed")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "internal server error" {
		t.Errorf("unexpected error message: %v", response["error"])
	}

	if int(response["code"].(float64)) != 500 {
		t.Errorf("unexpected code: %v", response["code"])
	}
}