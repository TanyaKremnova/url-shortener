package handlers_test

import (
	"bytes"
	"database/sql"
	"time"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"


	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/TanyaKremnova/url-shortener/internal/handlers"
)

// helper — creates a sqlx.DB backed by sqlmock
func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "postgres"), mock
}

// helper — creates a Gin context with a JSON body and a recorder
func newTestContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

// helper — parse response body into a map
func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("could not parse response body: %v\nbody was: %s", err, w.Body.String())
	}
	return result
}

// ── Register ──────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	db, mock := newMockDB(t)
	defer db.Close()

	// Expect INSERT ... RETURNING id
	rows := sqlmock.NewRows([]string{"id"}).AddRow("some-uuid-1234")
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()). // email, hash
		WillReturnRows(rows)

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "test@test.com",
		"password": "password123",
	})

	h.Register(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	body := parseBody(t, w)
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'data' key in response, got: %v", body)
	}
	if data["token"] == "" || data["token"] == nil {
		t.Error("expected token in response, got empty")
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "not-an-email",
		"password": "password123",
	})

	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_PasswordTooShort(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "test@test.com",
		"password": "short",
	})

	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{})

	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	// Simulate Postgres unique constraint violation (code 23505)
	pqErr := &pq.Error{Code: "23505"}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(pqErr)

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/register", map[string]string{
		"email":    "taken@test.com",
		"password": "password123",
	})

	h.Register(c)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d — body: %s", w.Code, w.Body.String())
	}

	body := parseBody(t, w)
	if body["error"] != "email already registered" {
		t.Errorf("unexpected error message: %v", body["error"])
	}
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	db, mock := newMockDB(t)
	defer db.Close()

	// Hash a known password so we can compare it in the handler
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
		AddRow("uuid-123", "test@test.com", string(hash), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE email`)).
		WithArgs("test@test.com").
		WillReturnRows(rows)

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "test@test.com",
		"password": "password123",
	})

	h.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	body := parseBody(t, w)
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'data' key, got: %v", body)
	}
	if data["token"] == "" || data["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "created_at"}).
		AddRow("uuid-123", "test@test.com", string(hash), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE email`)).
		WithArgs("test@test.com").
		WillReturnRows(rows)

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "test@test.com",
		"password": "wrongpassword",
	})

	h.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// Must not reveal which field was wrong — security requirement
	body := parseBody(t, w)
	if body["error"] != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got: %v", body["error"])
	}
}

func TestLogin_EmailNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE email`)).
		WithArgs("ghost@test.com").
		WillReturnError(sql.ErrNoRows)

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{
		"email":    "ghost@test.com",
		"password": "password123",
	})

	h.Login(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// Same message as wrong password — no enumeration
	body := parseBody(t, w)
	if body["error"] != "invalid credentials" {
		t.Errorf("expected 'invalid credentials', got: %v", body["error"])
	}
}

func TestLogin_MissingFields(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	h := handlers.NewAuthHandler(db)
	c, w := newTestContext(http.MethodPost, "/auth/login", map[string]string{})

	h.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}