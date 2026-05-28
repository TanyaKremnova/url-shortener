package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"database/sql"
	"time"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"

	"github.com/TanyaKremnova/url-shortener/internal/handlers"
)

type mockCache struct {
	data map[string]string
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string]string)}
}

func (m *mockCache) Get(_ context.Context, code string) (string, error) {
	url, ok := m.data["url:"+code]
	if !ok {
		return "", fmt.Errorf("cache miss")
	}
	return url, nil
}

func (m *mockCache) Set(_ context.Context, code, url string) error {
	m.data["url:"+code] = url
	return nil
}

func (m *mockCache) Delete(_ context.Context, code string) error {
	delete(m.data, "url:"+code)
	return nil
}

// ── Redirect ──────────────────────────────────────────────────────────────────

func TestRedirect_CacheHit(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	cache := newMockCache()
	cache.Set(context.Background(), "abc123", "https://www.github.com")

	// Expect the background click increment — use AnyArg since it runs in goroutine
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE urls`)).
		WithArgs("abc123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/abc123", nil)
	c.Params = gin.Params{{Key: "code", Value: "abc123"}}

	h := handlers.NewRedirectHandler(db, cache)
	h.Redirect(c)

	// Give goroutine time to run
	time.Sleep(10 * time.Millisecond)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", w.Code)
	}
	if w.Header().Get("Location") != "https://www.github.com" {
		t.Errorf("expected Location header to be https://www.github.com, got %s",
			w.Header().Get("Location"))
	}
}

func TestRedirect_CacheMiss_DBHit(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	cache := newMockCache() // empty cache

	rows := sqlmock.NewRows([]string{"original_url"}).
		AddRow("https://www.google.com")

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE urls`)).
		WithArgs("xyz789").
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/xyz789", nil)
	c.Params = gin.Params{{Key: "code", Value: "xyz789"}}

	h := handlers.NewRedirectHandler(db, cache)
	h.Redirect(c)

	if w.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", w.Code)
	}
	if w.Header().Get("Location") != "https://www.google.com" {
		t.Errorf("expected Location: https://www.google.com, got %s",
			w.Header().Get("Location"))
	}

	// Verify the URL is now in cache for next request
	cached, err := cache.Get(context.Background(), "xyz789")
	if err != nil {
		t.Error("expected URL to be cached after DB hit")
	}
	if cached != "https://www.google.com" {
		t.Errorf("expected cached URL to be https://www.google.com, got %s", cached)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	cache := newMockCache()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE urls`)).
		WithArgs("notexist").
		WillReturnError(sql.ErrNoRows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/notexist", nil)
	c.Params = gin.Params{{Key: "code", Value: "notexist"}}

	h := handlers.NewRedirectHandler(db, cache)
	h.Redirect(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	body := parseBody(t, w)
	if body["error"] != "short url not found" {
		t.Errorf("unexpected error: %v", body["error"])
	}
}