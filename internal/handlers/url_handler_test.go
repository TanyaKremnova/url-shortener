package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/TanyaKremnova/url-shortener/internal/auth"
	"github.com/TanyaKremnova/url-shortener/internal/handlers"
)

func TestCreateURL_Unauthorized(t *testing.T) {
	db, _, _ := sqlmock.New()
	h := handlers.NewURLHandler(sqlx.NewDb(db, "sqlmock"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/url", nil)
	c.Request = req

	h.CreateURL(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateURL_InvalidURL(t *testing.T) {
	db, _, _ := sqlmock.New()
	h := handlers.NewURLHandler(sqlx.NewDb(db, "sqlmock"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"original_url":"ftp://invalid"}`
	req := httptest.NewRequest("POST", "/url", bytes.NewBufferString(body))
	c.Request = req
	c.Set(auth.UserIDKey, "user-1")

	h.CreateURL(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateURL_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	h := handlers.NewURLHandler(sqlx.NewDb(db, "sqlmock"))

	os.Setenv("APP_BASE_URL", "http://short")

	mock.ExpectQuery("INSERT INTO urls").
		WillReturnRows(sqlmock.NewRows([]string{"short_code"}).AddRow("abc123"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"original_url":"https://google.com"}`
	req := httptest.NewRequest("POST", "/url", bytes.NewBufferString(body))
	c.Request = req
	c.Set(auth.UserIDKey, "user-1")

	h.CreateURL(c)

	require.Equal(t, http.StatusCreated, w.Code)
}