# Technical Notes

This file explains what I learned and why I made certain decisions for each part of the project.
It also includes what I would add with more time.

---

## 🎫 Ticket 1 — Docker + PostgreSQL

### Goal
Run a PostgreSQL database locally without installing it directly on the computer.

### What is Docker?
Docker runs software inside a "container" — an isolated box that has everything the software needs.
It's like a mini computer inside a computer. A container is similar with a virtual machines but much lighter and faster to start.

`docker-compose` is a tool that define containers in a file (`docker-compose.yml`) and start them with one command.

### Why PostgreSQL?
It is a reliable, open source relational database. It supports `UUID`, foreign keys, and transactions — all things this project needs. No strong reason to pick it over MySQL here; PostgreSQL is simply more common in Go projects.

### Database schema decisions

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ...
);
```

**UUID vs integer ID** — I used UUID instead of `SERIAL` (auto-increment integer).
- UUID: globally unique, cannot be guessed (`a3f8...` vs `1`, `2`, `3`)
- Integer: easier to read, but exposes how many users you have, and `GET /users/1` lets someone iterate through all users
- For a URL shortener where IDs might appear in URLs, UUID is safer
```bash
# in migrations/001_init.sql

id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

# UUID looks like: 550e8400-e29b-41d4-a716-446655440000
```

**`ON DELETE CASCADE`** on `urls.user_id` — if a user is deleted, all their URLs are deleted too. Otherwise we will get "orphan" rows with no owner in the database.

**Index on `short_code`**
```sql
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
```
Every redirect does `WHERE short_code = $1`. Without an index, PostgreSQL reads every row in the table to find the match. With an index it jumps directly to the right row.

### What I would add with more time
- Migration tool (e.g. `golang-migrate`) so schema changes are versioned and reversible
- Separate `events` table to log each click with timestamp and IP instead of just a counter

---

## 🎫 Ticket 2 — Database Connection + Config

### Goal
Connect the Go app to PostgreSQL and read configuration from environment variables.

### Why environment variables?
Never hardcode secrets (passwords, keys) in code. If the code is public on GitHub, anyone can read them.
Environment variables are set outside the code — in a `.env` file locally, or injected by the server in production.
This follows the **12-factor app** principle: configuration belongs in the environment, not in the code.

`.env` is in `.gitignore` — it never gets committed.
`.env.example` shows what variables are needed without real values — this one is committed

### What is a connection pool?
Instead of opening a new database connection for every request (slow), the app keeps a pool of open connections ready to reuse.
```go
db.SetMaxOpenConns(25)  // max 25 connections at the same time
db.SetMaxIdleConns(5)   // keep 5 open even when idle
```

---

## 🎫 Ticket 3 — HTTP Server + Router

### Goal
Start an HTTP server and define all routes.

### What is Gin?
Gin is a Go web framework. It handles incoming HTTP requests and routes them to the right function (handler).

### Why route groups?
```go
auth  := r.Group("/auth")   // /auth/register, /auth/login
urls  := r.Group("/urls")   // /urls/
admin := r.Group("/admin")  // /admin/urls/stats
```
Groups allow to apply middleware to multiple routes at once (e.g. auth check on all `/admin` routes). They also keep the code organised — easy to see which routes belong together.

### Why `/:code` must be last
Gin matches routes in the order they are defined. If `/:code` was first, it would catch `/health`, `/auth/register`, everything.
Defining it last means all specific routes are checked first.

### `/health` endpoint
Returns `{"status": "ok"}`. Used to check the server is running.
In production, load balancers and monitoring tools ping this endpoint automatically.

---

## 🎫 Ticket 4 — Auth: Register + Login

### Goal
Let users create an account and log in. Return a JWT token on success.

### What is bcrypt?
bcrypt is a hashing function for passwords. We never store the password itself — only the hash.
```
password: "hello123"
hash:     "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
```
When the user logs in, bcrypt compares the input with the stored hash.
It is designed to be slow on purpose — this makes brute-force attacks expensive.

### What is JWT?
JWT (JSON Web Token) is a signed string that proves who you are.

```
header.payload.signature
```

- **Header** — algorithm used to sign (`HS256`)
- **Payload** — data inside the token (`user_id`, `exp`)
- **Signature** — created with your secret key; proves the token was not changed

The server never stores the token. On every request, it just checks the signature.
This is called **stateless auth** — no session table in the database needed.

**Trade-off vs sessions:**
- JWT: stateless, scales easily, but you cannot invalidate a token before it expires
- Sessions: stored in DB, you can invalidate anytime, but every request needs a DB lookup

### Why same error message for wrong email and wrong password?
```go
c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
```
If you return "email not found" vs "wrong password" separately, an attacker can enumerate valid emails by trying many emails and watching for the different message. Same message for both = no information leak.

---

## 🎫 Ticket 5 — Auth Middleware

### Goal
Protect routes so only logged-in users can access them.

### What is middleware?
Middleware is a function that runs before the handler. It can check something and either continue (`c.Next()`) or stop the request (`c.Abort()`).

```
Request → Middleware (check token) → Handler (do the work)
                    ↓ if invalid
                 Return 401
```

### 401 vs 403
- `401 Unauthorized` — you are not logged in. "I don't know who you are."
- `403 Forbidden` — you are logged in, but not allowed. "I know who you are, but no."

In this project: no token = `401`, token from wrong user = `403`.

---

## 🎫 Ticket 6 — URL Shortening

### Goal
Accept a long URL, generate a short code, save it, return the short link.

### Short code generation
```go
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6
```
This is **base62** — 62 possible characters per position.
`62^6 = 56,800,235,584` — about 56 billion possible codes.

**`crypto/rand` vs `math/rand`**
- `math/rand` is deterministic — if you know the seed, you can predict all future values
- `crypto/rand` uses the operating system's secure random source — not predictable. In security contexts, always use the OS-level source.

### Why validate URLs server-side?
The `binding:"required"` tag checks the field is not empty.
`IsValidURL()` checks it actually starts with `http://` or `https://` and has a real host.
Never trust client input. Someone can send any string in the request body.

---

## 🎫 Ticket 7 — Redirect

### Goal
When someone visits the short URL, redirect them to the original URL and count the click.

### 301 vs 302
- `301 Moved Permanently` — browser caches this forever. Next time you visit the short URL, the browser goes directly to the destination without asking the server. **Click counter would never increase after the first visit.**
- `302 Found` — browser always asks the server. Counter increments every time.

Always use `302` for URL shorteners.

### Atomic update
```sql
UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1 RETURNING original_url
```
This does two things in one query: increment the counter AND return the original URL.

If we will use `SELECT` then `UPDATE` separately (two queries), two users clicking at the same time could both read `click_count = 5`, both write `6`, and one click would be lost (race conditions). One atomic query prevents this.

---

## 🎫 Ticket 8 — Stats Endpoint

### Goal
Return all URLs created by the logged-in user with their click counts.

### Why read `user_id` from the token, not the request?
```go
userID, _ := c.Get(auth.UserIDKey)  // from JWT — safe
```
If you read `user_id` from a query param or request body, any user could send `?user_id=someone_else` and see their data. The JWT is signed by the server — the client cannot change what is inside it.

### Empty array vs null
```go
if urls == nil {
    urls = []models.URLStats{}
}
```
When a user has no URLs, returning `{"urls": []}` is better than `{"urls": null}`.
Null forces every client to add a null check. Empty array is safe to iterate over directly.

---

## 🎫 Ticket 9 — Error Handling + Validation

### Goal
Every error response looks the same across the whole API.

### Consistent error shape
```json
{ "error": "message here", "code": 404 }
```
This makes the API predictable. A frontend or another service consuming this API always knows where to find the error message — it never has to guess the shape.

### Recovery middleware
```go
defer func() {
    if err := recover(); err != nil { ... }
}()
```
In Go, a `panic` is an unexpected crash (like a segfault in C). Without recovery middleware, one panic crashes the whole server for all users. With it, the panic is caught, logged, and the user gets a clean `500` response. The server keeps running.

---

## Feature (What I would add with more time)

### Rate limiting

Limit how many requests a user or IP can send in a period of time. This helps protect the API from spam, brute-force attacks, and abuse.

### Structured logging

Add logs in a structured JSON format to make debugging and monitoring easier in production systems.

### Dockerized API service

Run the Go API inside Docker together with PostgreSQL for a fully reproducible development environment.

### Integration tests

Add tests that verify complete application flows such as:

- register
- login
- create short URL
- redirect
- stats access

### Expiration dates for URLs

Allow short URLs to expire automatically after a chosen date or time.

### Custom aliases

Allow users to create custom short links such as: `/portfolio` instead of only random generated codes.

### Analytics event table

Store every redirect event separately with timestamp and metadata instead of only storing a total click counter.

---

## HTTP Status Codes Reference

| Code | Name | When used in this project |
|------|------|--------------------------|
| `200` | OK | Successful GET or login |
| `201` | Created | User registered, URL created |
| `302` | Found | Redirect to original URL |
| `400` | Bad Request | Invalid input (bad URL, missing field) |
| `401` | Unauthorized | Missing or invalid JWT token |
| `403` | Forbidden | Authenticated but accessing someone else's data |
| `404` | Not Found | Short code does not exist |
| `409` | Conflict | Email already registered |
| `500` | Internal Server Error | Unexpected server error |

---

## End-to-End Test

Full test of every feature from zero, using only `curl` in the terminal.

```bash
# ── 1. Health check ────────────────────────────────────────────────
curl -s http://localhost:8080/health
# {"status":"ok"}


# ── 2. Register two users ──────────────────────────────────────────
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@test.com", "password": "password123"}'
# {"code":201,"data":{"token":"eyJ..."}}

curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "bob@test.com", "password": "password123"}'
# {"code":201,"data":{"token":"eyJ..."}}


# ── 3. Try to register with the same email ─────────────────────────
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@test.com", "password": "password123"}'
# {"code":409,"error":"email already registered"}


# ── 4. Login as Alice, save token ──────────────────────────────────
ALICE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@test.com", "password": "password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Alice token: $ALICE"


# ── 5. Login as Bob, save token ────────────────────────────────────
BOB=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "bob@test.com", "password": "password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Bob token: $BOB"


# ── 6. Try to create URL without token ─────────────────────────────
curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://www.github.com"}'
# {"code":401,"error":"authorization header required"}


# ── 7. Try to create URL with invalid URL ──────────────────────────
curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE" \
  -d '{"original_url": "not-a-url"}'

# {"code":400,"error":"invalid url: must start with http:// or https://"}


# ── 8. Alice creates two URLs ──────────────────────────────────────
CODE1=$(curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE" \
  -d '{"original_url": "https://www.github.com"}' \
  | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)

echo "Alice code 1: $CODE1"

CODE2=$(curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ALICE" \
  -d '{"original_url": "https://www.google.com"}' \
  | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)

echo "Alice code 2: $CODE2"


# ── 9. Bob creates one URL ─────────────────────────────────────────
CODE3=$(curl -s -X POST http://localhost:8080/urls/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $BOB" \
  -d '{"original_url": "https://www.codam.nl"}' \
  | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)

echo "Bob code: $CODE3"


# ── 10. Visit Alice's first URL 3 times ───────────────────────────
curl -sL http://localhost:8080/$CODE1 -o /dev/null
curl -sL http://localhost:8080/$CODE1 -o /dev/null
curl -sL http://localhost:8080/$CODE1 -o /dev/null

# Visit Alice's second URL 1 time
curl -sL http://localhost:8080/$CODE2 -o /dev/null

# Visit Bob's URL 5 times
curl -sL http://localhost:8080/$CODE3 -o /dev/null
curl -sL http://localhost:8080/$CODE3 -o /dev/null
curl -sL http://localhost:8080/$CODE3 -o /dev/null
curl -sL http://localhost:8080/$CODE3 -o /dev/null
curl -sL http://localhost:8080/$CODE3 -o /dev/null


# ── 11. Try a code that does not exist ─────────────────────────────
curl -s http://localhost:8080/doesnotexist
# {"code":404,"error":"short url not found"}


# ── 12. Alice checks her stats ─────────────────────────────────────
curl -s http://localhost:8080/admin/urls/stats \
  -H "Authorization: Bearer $ALICE"

# Should show CODE1 with click_count=3 and CODE2 with click_count=1
# Bob's URL should NOT appear here


# ── 13. Bob checks his stats ───────────────────────────────────────
curl -s http://localhost:8080/admin/urls/stats \
  -H "Authorization: Bearer $BOB"

# Should show CODE3 with click_count=5
# Alice's URLs should NOT appear here


# ── 14. Stats without token ────────────────────────────────────────
curl -s http://localhost:8080/admin/urls/stats
# {"code":401,"error":"authorization header required"}
```