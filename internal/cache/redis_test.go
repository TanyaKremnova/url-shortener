package cache

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func setupTestCache() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Cache{
		client: client,
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := setupTestCache()
	ctx := context.Background()

	code := "abc123"
	url := "https://github.com"

	err := cache.Set(ctx, code, url)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	result, err := cache.Get(ctx, code)
	if err != nil {
		t.Fatalf("failed to get cache: %v", err)
	}

	if result != url {
		t.Errorf("expected %s, got %s", url, result)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := setupTestCache()
	ctx := context.Background()

	code := "delete123"
	url := "https://google.com"

	err := cache.Set(ctx, code, url)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	err = cache.Delete(ctx, code)
	if err != nil {
		t.Fatalf("failed to delete cache: %v", err)
	}

	_, err = cache.Get(ctx, code)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}