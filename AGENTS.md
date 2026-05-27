# Gater

Event ticketing API. Go, PostgreSQL, Redis.

## Quick start

```sh
just dev       # docker compose up + air (hot reload)
just migrate   # run goose migrations
go build -o bin/server ./cmd/server
```

## Project layout

```
cmd/server/        -- main, HTTP handlers, routing, middleware (package main)
cmd/migrate/       -- goose migration runner
internal/
  config/          -- env var loading (godotenv)
  db/              -- pgxpool creation
  cache/           -- redis.Client creation
  store/           -- data access (raw SQL via pgx), Store{Users,Sessions,Verifications}
  auth/            -- argon2id hashing, SHA-256 token gen/verify
  json/            -- JSON HTTP helpers (1MB limit, DisallowUnknownFields)
  mailer/          -- Mailer interface + Resend implementation
  validator/       -- go-playground/validator wrapper
```

## Key conventions

- **Router:** `go-chi/chi/v5`. Global middleware stack in `mount()`: CleanPath, StripSlashes, RequestID, RealIP, Logger, Recoverer, CORS, 60s timeout.
- **No DI framework** -- everything manually wired in `main.go` into the `application` struct.
- **Handlers** live in `cmd/server/`, **infra** in `internal/`.
- **Password:** argon2id (64MB mem, 3 iter, 4 threads).
- **Session tokens:** 32-byte random hex, SHA-256 hashed, constant-time compare. Cookie: `gater_auth_session`, HttpOnly/Secure/Lax, 30d.
- **Background email** uses `context.Background()` (not request ctx), errors only logged.
- **No database transactions** -- store methods are individual queries, no rollback.
- **`requireAuth` middleware is a no-op** -- context injection is commented out.
- **Handler stubs** exist for events, tiers, purchases, tickets, waitlist, check-in, analytics -- all empty `package main` files.

## Routes (all under `/api`)

```
GET  /api/health
POST /api/auth/register, login, verify-email, resend-verification, forgot-password, reset-password
GET  /api/auth/google, google/callback
POST /api/auth/logout (protected)
GET  /api/auth/me (protected)
```

## Testing

No tests exist. All `*_test.go` files would go next to the code they test.

## Migrations

Run with `just migrate` (goose, postgres dialect, embedded SQL in `cmd/migrate/migrations/`).
10 migrations cover users, sessions, oauth_accounts, verifications, events, ticket_tiers, purchases, tickets, waitlist_entries. Each `updated_at` column gets an auto-trigger.

## Environment

`.env` is gitignored. Copy `.env.example` for defaults. Port 8080, PG on 5435, Redis on 6380.

## Known quirks

- `go.mod` module path: `github.com/chiagxziem/gater` (typo: `chiagxziem` not `chiagoziem`).
- No CSRF protection (CORS allows `X-CSRF-Token` but it's never validated).
- No rate limiting beyond verification-resend's ad-hoc check.
