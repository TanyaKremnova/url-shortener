package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TanyaKremnova/url-shortener/internal/config"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()

	old, existed := os.LookupEnv(key)

	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})

	_ = os.Setenv(key, value)
}

func TestLoad_ConfigFromEnv(t *testing.T) {
	setEnv(t, "DATABASE_URL", "postgres://db")
	setEnv(t, "APP_PORT", "8080")
	setEnv(t, "JWT_SECRET", "secret")
	setEnv(t, "APP_BASE_URL", "http://localhost")
	setEnv(t, "CACHE_URL", "redis://localhost")

	cfg := config.Load()

	require.Equal(t, "postgres://db", cfg.DatabaseURL)
	require.Equal(t, "8080", cfg.AppPort)
	require.Equal(t, "secret", cfg.JWTSecret)
	require.Equal(t, "http://localhost", cfg.AppBaseURL)
	require.Equal(t, "redis://localhost", cfg.CacheURL)
}

func TestLoad_EmptyEnv(t *testing.T) {
	// isolate environment
	keys := []string{
		"DATABASE_URL",
		"APP_PORT",
		"JWT_SECRET",
		"APP_BASE_URL",
		"CACHE_URL",
	}

	for _, k := range keys {
		_ = os.Unsetenv(k)
	}

	cfg := config.Load()

	require.Equal(t, "", cfg.DatabaseURL)
	require.Equal(t, "", cfg.AppPort)
	require.Equal(t, "", cfg.JWTSecret)
	require.Equal(t, "", cfg.AppBaseURL)
	require.Equal(t, "", cfg.CacheURL)
}