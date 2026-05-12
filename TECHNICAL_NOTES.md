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

# 🎫 Ticket 2 — Database Connection + Config

# 🎫 Ticket 3 — HTTP Server + Router

# 🎫 Ticket 4 — Auth: Register + Login

# 🎫 Ticket 5 — Auth Middleware

# 🎫 Ticket 6 — URL Shortening

# 🎫 Ticket 7 — Redirect

# 🎫 Ticket 8 — Stats Endpoint

# 🎫 Ticket 9 — Error Handling + Validation

# 🎫 Ticket 10 — README + Docs + Final Polish