# Gater

Event ticketing API. Go 1.26, PostgreSQL, Redis.

## Quick start

```sh
just dev       # docker compose up -d + air (hot reload)
just d         # alias for dev
just migrate   # runs goose migrations (go run cmd/migrate/main.go)
just db-up     # docker compose up -d (PG:5435, Redis:6380)
just db-down   # docker compose down
just db-delete # docker compose down -v
go build -o bin/server ./cmd/server
```

No tests, no linter, no formatter config.

## Architecture

`cmd/server/` — `package main`, HTTP handlers, chi routing, middleware.  
`internal/` — `config/` (godotenv), `db/` (pgxpool), `store/` (raw SQL via pgx, 5s per-query timeout), `auth/` (argon2id, SHA-256 tokens), `jsonutil/`, `validator/` (go-playground), `mailer/` (Resend).  
`cmd/migrate/` — goose runner with embedded SQL.  
`internal/cache/redis.go` exists but is **not wired into the app**.

Handlers manually wired into `application` struct in `main.go` — no DI framework.

## Key conventions

- **Auth:** `requireAuth` middleware checks `Authorization: Bearer <token>` first, then falls back to the `gater_auth_session` cookie for browser clients. Token = 32-byte random → hex → SHA-256 → store hash. Session create retries up to 3× on hash collision.
- **Cookie** `gater_auth_session` set on login (HttpOnly, Lax, 30d, Secure only in production, `SameSite=Lax`). CORS `AllowCredentials: true` lets browsers send it cross-origin.
- **JSON response:** Success `{"data": ...}` via `WriteData`, errors `{"errors": [...]}` via `WriteError`. Exception: health check uses bare `Write` → `{"status":"OK"}`.
- **Password** `json:"-"` — never serialized to JSON. `internal/store/` uses raw SQL, no transactions.
- **Background email** uses `context.Background()`, errors only logged.
- **Ad-hoc rate limiting** on verification-resend and forgot-password: 5 per hour, 1 min cooldown, checked via `Verifications.CountSince`.

## Routes (`/api`)

```
GET  /api/health
POST /api/auth/register, login, verify-email, resend-verification, forgot-password, reset-password
GET  /api/auth/google, google/callback
POST /api/auth/logout, become-organizer  (protected)
GET  /api/auth/me                        (protected)
```

## Implemented vs stubs

| File                                                                                  | Status      |
| ------------------------------------------------------------------------------------- | ----------- |
| `auth.go`, `users.go`, `health.go`, `middleware.go`                                   | Implemented |
| `events.go`, `tiers.go`, `purchases.go`, `check-in.go`, `waitlist.go`, `analytics.go` | Empty stubs |

## Requests

`requests/` directory contains a [Bruno](https://docs.usebruno.com) API collection (`opencollection.yml`) with request examples for auth endpoints.

## Known quirks

- `internal/cache/redis.go` imports redis client but nothing in the app uses it.
