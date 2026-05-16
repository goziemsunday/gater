# AGENTS.md

## Dev Commands

```bash
# Start all services (Docker + overmind)
just dev

# Stop services
just down

# Server only: hot-reload dev
cd server && just dev

# Server only: build
cd server && just build
```

## Architecture

- **Monorepo**: `server/` (Go), `web/` (TanStack Start - empty/scaffolding)
- **Backend entrypoint**: `server/cmd/server/main.go`
- **Go module**: `github.com/chiagxziem/snipper`

## Infrastructure

- PostgreSQL: port **5435** (not 5432)
- Redis: port **6380** (not 6379)
- Docker Compose: `docker-compose.yaml` at repo root

## Current State

Early scaffolding. Server has minimal stub code (`fmt.Println("hello world!")`). No frontend, no auth, no URL shortener logic yet.

## Dependencies

Go dependencies managed via `server/go.mod`. No frontend dependencies yet.
