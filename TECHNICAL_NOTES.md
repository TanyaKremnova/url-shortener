# 🎫 Ticket 0 — Project Bootstrap

```bash
go mod init github.com/yourname/url-shortener
```

### go.mod

The go.mod file defines:
- module name
- Go version
- dependencies

### go.sum - The Security Lock
While go.mod says which version you need, **go.sum** ensures the code you download is exactly what you expect and hasn't been tampered with.

It contains:
- A list of cryptographic checksums (hashes) for the source code and go.mod files of every dependency.

It exists:
- **Integrity**: It prevents supply chain attacks. If a hacker modifies a library you depend on, the hash won't match your go.sum file, and Go will refuse to build.
- **Reproducibility**: It guarantees that every developer on your team and your CI/CD pipeline are using the exact same "bits" of code, even if a library author tries to change a published version later.

***

# 🎫 Ticket 1 — Docker + PostgreSQL

```bash
services:
  postgres:
    image: postgres:17 #official PostgreSQL image version 17
    container_name: url_shortener_db

    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}

    ports:
      - "${POSTGRES_PORT}:5432" #5432:5432 -> local_machine:container

    volumes:
      - postgres_data:/var/lib/postgresql/data #stores DB data permanently
      - ./migrations:/docker-entrypoint-initdb.d #mounts the SQL files into container
	  # Postgres image automatically executes .sql files there on first startup

volumes:
  postgres_data:
```
***
## UUID vs SERIAL
### UUID Advantages:

- globally unique
- harder to guess
- safer for public APIs
- better for distributed systems

```bash
# in migrations/001_init.sql

id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

# UUID looks like: 550e8400-e29b-41d4-a716-446655440000
```

### SERIAL
```bash
id SERIAL PRIMARY KEY

# Simple incrementing integers
```

Good:
- smaller
- faster
- simpler

Bad:
- predictable
- easier enumeration attacks

Example:
- attacker guesses /users/1
- /users/2

## Why UUID Instead of Integer IDs

I chose UUIDs because they are harder to guess and safer for public-facing APIs.
Unlike incremental integer IDs, UUIDs reduce predictability and are useful in distributed systems where unique identifiers must be generated independently.

The tradeoff is larger storage size and slightly slower indexing compared to SERIAL integer IDs.

***

```bash
user_id UUID REFERENCES users(id),

# A foreign key means -> Every url.user_id must belong to an existing user

# Without this:
# - orphan data possible
# - broken relationships possible
```

### click_count
For MVP click_count lives on the urls table. Simple and efficient way:

```bash
CREATE TABLE urls (
	...
	click_count INT DEFAULT 0,
	...
```

Real systems
```bash
click_events
- url_id
- timestamp
- ip
- user_agent
```

***

## Useful Docker commands
```bash
docker-compose up -d          # start in background
docker-compose down           # stop containers
docker-compose down -v        # stop AND delete volume (wipes DB data)
docker logs url_shortener_db  # see postgres logs
docker ps                     # see running containers

# If SQL changed and need to re-run it:
docker-compose down -v        # wipe everything
docker-compose up -d          # start fresh — SQL runs again

# ⚠️ down -v deletes all data. Use it if migration SQL was changed.
```

## Start Postgres
```bash
docker-compose up -d
# ✅ Container url_shortener_db  Started
docker ps
docker logs url_shortener_db
# database system is ready to accept connections
```

```bash
# Connect to Postgres inside the container:
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
***

# 🎫 Ticket 2 — Database Connection + Config

## Goal

The goal of this ticket was to connect the Go application to PostgreSQL while keeping configuration separated from application logic.

The application now:

- loads configuration from environment variables
- connects to PostgreSQL using sqlx
- uses Dockerized PostgreSQL
- initializes database connection pooling

### Why Use Environment Variables

Database credentials and application configuration should not be hardcoded inside the source code.

Using environment variables provides:

- safer configuration management
- easier deployment to different environments
- flexibility between local development and production

Example:

- local development may use localhost
- production may use managed cloud databases

The **.env** file is ignored by Git to avoid committing secrets.

The **.env.example** file documents which variables are required without exposing real credentials.

### Why Use sqlx Instead of database/sql

Go's standard database/sql package is low-level and verbose.

**sqlx** extends it with:

- easier query helpers
- automatic struct scanning
- cleaner query handling

### Why Docker Is Useful Here

Docker provides a reproducible development environment.

Instead of manually installing PostgreSQL locally, the database runs inside a container with:

- fixed version
- isolated environment
- persistent storage through Docker volumes

This makes onboarding and setup more predictable.

***

# 🎫 Ticket 3 — HTTP Server + Router

***

# 🎫 Ticket 4 — Auth: Register + Login

## Goal

Implement user authentication using:

- user registration
- password hashing with bcrypt
- JWT-based login authentication

The API now allows users to:

- create accounts securely
- authenticate using email/password
- receive signed JWT access tokens

Tests:
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

```bash
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"supersecret"
}'
# HTTP 201 Created
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"supersecret"
}'
# HTTP 409 Conflict
```

```bash
curl -X POST http://localhost:8080/auth/register \
-H "Content-Type: application/json" \
-d '{}'
# HTTP 400 Bad Request
```

```bash
curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"supersecret"
}'
# HTTP 200 OK
```

```bash
curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email":"test@example.com",
  "password":"wrongpassword"
}'
# HTTP 401 Unauthorized
```

## Why Passwords Must Be Hashed

Passwords should never be stored in plain text.

If a database leak happens and passwords are unhashed:

- every user account becomes compromised immediately

bcrypt hashes passwords before storing them.

bcrypt also automatically includes:

- salt
- configurable computational cost

This makes brute-force attacks significantly harder.

## What Is Salt

A salt is random data added to a password before hashing.

Without salt:

- identical passwords produce identical hashes

With salt:

- same passwords generate different hashes

This protects against rainbow table attacks.

bcrypt automatically handles salt generation internally.

## Why bcrypt.CompareHashAndPassword Is Important

Password comparison should be timing-safe.

Simple string comparison can leak timing information and theoretically help attackers.

bcrypt.CompareHashAndPassword performs secure hash verification designed for authentication systems.

## What Is JWT

JWT (JSON Web Token) is a signed authentication token.

A JWT contains:

- header
- payload
- signature

Example payload:

- user ID
- expiration timestamp (exp)

The signature prevents token tampering.

## Stateless Authentication

JWT authentication is stateless.

The server does not store session state in memory or database tables.

Instead:

- client stores token
- token is sent with requests
- server validates signature

Benefits:

- simpler horizontal scaling
- no session storage needed
- useful for APIs and microservices

Tradeoffs:

- harder token revocation
- logout is less straightforward
- token expiration strategy becomes important

## HTTP Status Codes Used
### 200 OK
### 201 Created

Used when a new user account is successfully created.

### 400 Bad Request

Used when request input is invalid or malformed.

### 401 Unauthorized

Used when login credentials are invalid.

### 409 Conflict

Used when trying to register an email that already exists.

***

# 🎫 Ticket 5 — Auth Middleware

## What Is Middleware

Middleware is a function that executes before the main route handler.

It is commonly used for:

- authentication
- logging
- rate limiting
- request validation

Middleware can:

- stop the request
- modify the request context
- pass control to the next handler

## JWT Authentication Middleware

The authentication middleware:

- reads the Authorization header
- extracts the Bearer token
- validates the JWT signature and expiration
- rejects invalid or missing tokens

If validation succeeds:

- the authenticated user_id is stored in the request context

Downstream handlers can then access the authenticated user without re-validating the token.

## 401 Unauthorized

Used when:

- token is missing
- token is invalid
- authentication failed

Meaning:

"You are not authenticated."

## 403 Forbidden

Used when:

- user is authenticated
- but does not have permission

Meaning:

"You are authenticated, but not allowed to access this resource."

Tests:
```bash
curl -X POST http://localhost:8080/urls/
# HTTP 401
```

```bash
# First get a token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "password123"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Then use it
curl -X POST http://localhost:8080/urls/ \
  -H "Authorization: Bearer $TOKEN"

# {"message":"not implemented"}  ← middleware passed! ✅
# HTTP 501
```

***

# 🎫 Ticket 6 — URL Shortening

***

# 🎫 Ticket 7 — Redirect

***

# 🎫 Ticket 8 — Stats Endpoint

***

# 🎫 Ticket 9 — Error Handling + Validation

***

# 🎫 Ticket 10 — README + Docs + Final Polish