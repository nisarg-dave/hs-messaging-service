# AGENTS.md

Guide for AI coding agents (and humans) working in `hs-messaging-service`.

## What this is

A small Go HTTP service for sending and reading messages. Stack:

- **Language:** Go (see `go.mod` for version)
- **HTTP:** [Echo v5](https://github.com/labstack/echo)
- **DB:** Postgres via [GORM](https://gorm.io)
- **Local infra:** Postgres in `docker-compose.yml`

## Repo layout

```
cmd/api/              # main entrypoint
internal/
  api/
    handlers/         # HTTP handlers (Echo)
    routes/           # Route registration
  service/            # Business logic
  repository/postgres # GORM data access
  domain/             # Domain models (GORM tags)
  config/             # Env-based config loader
docker-compose.yml    # Local Postgres
.cursor/rules/        # Cursor-specific agent rules (authoritative for Cursor)
```

Request flow: `routes -> handlers -> service -> repository -> domain`.
Keep each layer's responsibilities narrow; don't let handlers talk to GORM directly and don't let repositories return Echo types.

## How agents should work here

1. **Follow `.cursor/rules/`** — those are the binding conventions for Cursor. This file summarizes; the rules are authoritative.
2. **Small, layered changes.** New endpoints usually mean: domain model (if needed) -> repository method -> service method -> handler -> route registration in `cmd/api/main.go`.
3. **Never commit secrets.** `.env` is gitignored; don't check in credentials or real connection strings.

## Common commands

```bash
# Start local Postgres
docker compose up -d

# Run the API (reads .env via godotenv)
go run ./cmd/api

# Build
go build ./...

# Format & vet
gofmt -w .
go vet ./...

# Tests (when present)
go test ./...
```

## Environment

`.env` is loaded at startup. Expected keys:

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- `SERVER_PORT`

## Do / don't

- Do keep handlers thin: parse input, call service, return JSON.
- Do put validation and business rules in `service/`.
- Do add a GORM `AutoMigrate` entry in `repository/postgres/connection.go` when you add a new domain model.
- Don't log or return raw internal errors to clients without mapping to an HTTP status.
- Don't introduce new frameworks without discussion.
