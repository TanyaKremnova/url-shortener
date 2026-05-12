# 🎫 Ticket 0 — Project Bootstrap

```bash
go mod init github.com/yourname/url-shortener
```

### go.mod — The Manifest
go.mod - as the blueprint. It is located at the root of your project and defines the module's path and its dependency requirements

It contains:
- **Module Path**: The unique name of a module (often a URL like ://github.com).
- **Go Version**: The minimum version of Go required for the module.
- **Require Directives**: A list of direct and some indirect dependencies along with their specific semantic versions (e.g., v1.2.3).

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


# 🎫 Ticket 2 — Database Connection + Config

***

# 🎫 Ticket 3 — HTTP Server + Router

***

# 🎫 Ticket 4 — Auth: Register + Login

***

# 🎫 Ticket 5 — Auth Middleware

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