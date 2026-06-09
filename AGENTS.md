# Gater

Event ticketing API. Go, PostgreSQL, Redis.

## Quick start

```sh
just dev       # docker compose up -d + air (hot reload)
just d         # alias for dev
just migrate   # run goose migrations
go build -o bin/server ./cmd/server
```

`just dev` automatically starts PG/Redis containers before launching air.

## Project layout

```
cmd/server/        -- main, HTTP handlers, routing, middleware (package main)
cmd/migrate/       -- goose migration runner
internal/
  config/          -- env var loading (godotenv)
  db/              -- pgxpool creation
  cache/           -- redis.Client creation (exists but not wired into app yet)
  store/           -- data access (raw SQL via pgx), Store{Users,Sessions,Verifications,OAuthAccounts}
  auth/            -- argon2id hashing, SHA-256 token gen/verify, OAuth state gen
  jsonutil/        -- JSON HTTP helpers (1MB limit, DisallowUnknownFields), envelope responses
  mailer/          -- Mailer interface + Resend implementation
  validator/       -- go-playground/validator wrapper
```

## Key conventions

- **Router:** `go-chi/chi/v5`. Global middleware in `mount()`: CleanPath, StripSlashes, RequestID, RealIP, Logger, Recoverer, CORS, 60s timeout, then custom `injectLogging`.
- **No DI framework** -- everything manually wired in `main.go` into the `application` struct.
- **Handlers** live in `cmd/server/`, **infra** in `internal/`.
- **Password:** argon2id (64MB mem, 3 iter, 4 threads). See `internal/auth/password.go`.
- **Session tokens:** 32-byte random hex, SHA-256 hashed, constant-time compare. Cookie: `gater_auth_session`, HttpOnly, Secure (prod only), Lax, 30d.
- **Background email** uses `context.Background()` (not request ctx), errors only logged.
- **No database transactions** -- store methods are individual queries, no rollback.
- **JSON response envelope:** Success `{"data": ...}`, errors `{"errors": [...]}`. See `internal/jsonutil/json.go`.
- **Handler stubs** exist for events, tiers, purchases, waitlist, check-in, analytics -- all empty `package main` files.
- **`requireAuth`** is fully implemented: extracts Bearer token, hashes, looks up session, injects user + session into context.

## Routes (all under `/api`)

```
GET  /api/health
POST /api/auth/register, login, verify-email, resend-verification, forgot-password, reset-password
GET  /api/auth/google, google/callback
POST /api/auth/logout, become-organizer  (protected)
GET  /api/auth/me                        (protected)
```

## Testing

No tests exist. All `*_test.go` files would go next to the code they test.

## Migrations

Run with `just migrate` (goose, postgres dialect, embedded SQL in `cmd/migrate/migrations/`).
10 migrations cover users, sessions, oauth_accounts, verifications, events, ticket_tiers, purchases, tickets, waitlist_entries. Each `updated_at` column gets an auto-trigger.

## Environment

`.env` is gitignored. Copy `.env.example` for defaults. Port 8080, PG on 5435, Redis on 6380.
All config fields are required (except PORT which defaults to 8080).

## Docker commands

```sh
just db-up      # docker compose up -d (PG + Redis)
just db-down    # docker compose down
just db-delete  # docker compose down -v (deletes volumes)
```

## Known quirks

- `go.mod` module path: `github.com/chiagxziem/gater` (typo: `chiagxziem` not `chiagoziem`).
- No CSRF protection (CORS allows `X-CSRF-Token` but it's never validated).
- No rate limiting beyond verification-resend's ad-hoc check (5 per hour, 1 min cooldown).
- `internal/cache/redis.go` exists but isn't imported or used anywhere in the app yet.
- No linter/formatter config in repo -- no `.golangci-lint.yml` or equivalent.
- `.air.toml` watches all `.go` files, excludes `assets/`, `tmp/`, `vendor/`, `testdata/`.
