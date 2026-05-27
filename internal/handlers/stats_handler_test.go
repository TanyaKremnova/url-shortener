package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/TanyaKremnova/url-shortener/internal/auth"
	"github.com/TanyaKremnova/url-shortener/internal/handlers"
)

func TestGetStats_Unauthorized(t *testing.T) {
	db, _, _ := sqlmock.New()
	h := handlers.NewStatsHandler(sqlx.NewDb(db, "sqlmock"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/stats", nil)
	c.Request = req

	h.GetStats(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetStats_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	h := handlers.NewStatsHandler(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT").
		WillReturnError(sql.ErrConnDone)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/stats", nil)
	c.Request = req
	c.Set(auth.UserIDKey, "user-1")

	h.GetStats(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}