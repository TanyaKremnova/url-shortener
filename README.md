# URL Shortener

Backend URL shortener service written in Go.

# How to run and use

```bash
docker-compose up -d
go run ./cmd/api/main.go
```

```bash
# 1. Register
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "email@test.com", "password": "password123"}'

# 2. Login + grab token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "email@test.com", "password": "password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 3. Create URL
CODE=$(curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url": "https://www.github.com"}' \
  | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)

echo "Short code: $CODE"

# 4. Click the short URL 3 times
curl -sL http://localhost:8080/$CODE -o /dev/null
curl -sL http://localhost:8080/$CODE -o /dev/null
curl -sL http://localhost:8080/$CODE -o /dev/null

# 5. Check stats
curl -s http://localhost:8080/admin/urls/stats \
  -H "Authorization: Bearer $TOKEN"

# Expected: click_count = 3

```

### Errors
```bash
# 400 — bad request
curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url": "not-a-url"}'
# {"error":"invalid url: must start with http:// or https://","code":400}


# 401 — no token
curl -s http://localhost:8080/admin/urls/stats
# {"error":"authorization header required","code":401}


# 404 — bad code
curl -s http://localhost:8080/doesnotexist
# {"error":"short url not found","code":404}


# 409 — duplicate email
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "final@test.com", "password": "password123"}'
# {"error":"email already registered","code":409}
```