# AGENTS.md

## Dev Commands

```bash
# Start all services (Docker + Overmind)
just dev

# Stop services
just down

# Server dev (hot reload via Air)
cd server && just dev

# Build server binary → bin/server
cd server && just build

# Run database migrations (Goose, embedded SQL)
cd server && just migrate
```

## Architecture

- **Monorepo** with `server/` (Go). No frontend yet.
- **Backend entrypoint**: `server/cmd/server/main.go`
- **Go module**: `github.com/chiagxziem/snipper`
- **Router**: `go-chi/chi/v5` with routes defined in `api.go:mount()`
- **Store layer**: `server/internal/store/` wraps `pgx` pool
- **Migrations**: Goose v3 with embedded SQL at `server/cmd/migrate/migrations/*.sql`
- **Bruno API client**: `bruno-reqs/` — use Bruno app to hit `http://localhost:8080`

## Infrastructure

- PostgreSQL: port **5435** (not 5432)
- Redis: port **6380** (not 6379)
- Docker Compose: `docker-compose.yaml` at repo root
- Overmind manages processes via `Procfile`
- `.env` required: `PORT`, `DATABASE_URL`, `REDIS_URL`, `CORS_ALLOWED_ORIGIN`, `RESEND_API_KEY`, `RESEND_DOMAIN` (see `server/.env.example`)

## Current State

Early development. Auth register endpoint works; login, logout, URL CRUD, analytics, redirect, Google OAuth, email verification, password reset all wired but return 501. No tests. No frontend.

## VCS

Repo uses Jujutsu (`jj`) alongside git.
