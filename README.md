# URL Shortener

A backend REST API that shortens URLs. Built with Go, PostgreSQL, and Docker.

Registered users can create short links, share them, and see how many times each link was clicked.

This project was developed as part of a backend mentorship/take-home style exercise focused on API design, authentication, PostgreSQL, and backend architecture.

---

## Tech Stack

- **Go** — main language
- **Gin** — HTTP framework
- **PostgreSQL** — database
- **Docker + docker-compose** — runs the database locally
- **JWT** — authentication

---

## Requirements

- [Go 1.21+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

---

## How to Run

**1. Clone the repo**
```bash
git clone https://github.com/TanyaKremnova/url-shortener.git
cd url-shortener
```

**2. Create your `.env` file**
```bash
cp .env.example .env
```

Fill in `.env`:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable
APP_PORT=8080
JWT_SECRET=your_secret_key_here
APP_BASE_URL=http://localhost:8080
```

**3. Start the database**
```bash
docker-compose up -d
```

**4. Run the app**
```bash
go run ./cmd/api/main.go
```

The server is now running at `http://localhost:8080`.

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | No | Health check |
| `POST` | `/auth/register` | No | Create account |
| `POST` | `/auth/login` | No | Login, get JWT token |
| `POST` | `/urls/` | Yes | Create short URL |
| `GET` | `/:code` | No | Redirect to original URL |
| `GET` | `/admin/urls/stats` | Yes | See your URLs and click counts |

Protected routes require header: `Authorization: Bearer <token>`

---

## Quick Example

```bash
# Register
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "password123"}'

# Login and save token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Create short URL
curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url": "https://www.github.com"}'

# Visit short URL in browser
# http://localhost:8080/<your_code>

# Check stats
curl -s http://localhost:8080/admin/urls/stats \
  -H "Authorization: Bearer $TOKEN"
```

---

## Error Responses

All errors follow the same shape:
```json
{ "error": "message here", "code": 400 }
```

| Code | Meaning |
|------|---------|
| `400` | Bad request — invalid input |
| `401` | Unauthorized — missing or invalid token |
| `404` | Not found — short code does not exist |
| `409` | Conflict — email already registered |
| `500` | Server error |

---

### Error tests
```bash
# 400 — bad request
curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"original_url": "not-a-url"}'
# {"code":400,"error":"invalid url: must start with http:// or https://"}


# 401 — no token
curl -s http://localhost:8080/admin/urls/stats
# {"code":401,"error":"authorization header required"}


# 404 — bad code
curl -s http://localhost:8080/doesnotexist
# {"code":404,"error":"short url not found"}


# 409 — duplicate email
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "final@test.com", "password": "password123"}'
# {"code":409,"error":"email already registered"}
```

---

## Useful Commands

```bash
# Stop database
docker-compose down

# Wipe database and start fresh (deletes all data)
docker-compose down -v && docker-compose up -d
# ⚠️ down -v deletes all data. Use it if migration SQL was changed.

# Connect to database directly
docker exec -it url_shortener_db psql -U postgres -d url_shortener

SELECT * FROM users;

# Inside the Postgres shell. Run:
\dt

# Check the columns look right:
\d users
\d urls

# Exit the shell:
\q
```