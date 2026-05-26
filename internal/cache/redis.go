package cache

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
}

func New(addr string) *Cache {
    client := redis.NewClient(&redis.Options{Addr: addr})

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to cache: %v", err)
    }

    log.Println("✅ Connected to cache")
    return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, code, url string) error {
    return c.client.Set(ctx, "url:"+code, url, 24*time.Hour).Err()
}

func (c *Cache) Get(ctx context.Context, code string) (string, error) {
    return c.client.Get(ctx, "url:"+code).Result()
}

func (c *Cache) Delete(ctx context.Context, code string) error {
    return c.client.Del(ctx, "url:"+code).Err()
}