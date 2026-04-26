# hs-messaging-service

HTTP API for creating messages and marking them as read. Built with **Go**, **[Echo v5](https://github.com/labstack/echo)**, **GORM**, and **PostgreSQL**.

## Requirements

- [Go](https://go.dev/dl/) 1.26.2+ (see `go.mod`)
- [Docker](https://docs.docker.com/get-docker/) (for local Postgres)

## Quick start

### 1. Environment

Create a `.env` file in the project root (it is gitignored). `docker compose` and the app both expect these variables:

```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
SERVER_PORT=8080
```

Use the same database name, user, and password in `.env` as you pass to Postgres via Compose so the API can connect after the container is up.

### 2. Database

```bash
docker compose up -d
```

Wait until Postgres is healthy (`docker compose ps`).

### 3. Run the API

```bash
go run ./cmd/api
```

The server listens on the port in `SERVER_PORT` (for example `http://localhost:8080`).

## API

Base path for message routes: `/messages`.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/messages` | Create a message |
| `PATCH` | `/messages/:id/read` | Set `isRead` to true for the given message ID |

### Create message

**Request:** `POST /messages` with JSON body:

```json
{
  "senderId": "550e8400-e29b-41d4-a716-446655440000",
  "recipientId": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "content": "Hello",
  "jobId": "optional-uuid-or-omit"
}
```

`id`, `isRead`, `createdAt`, and `updatedAt` are managed by the database; you do not need to send them.

**Response:** `201 Created` with the persisted message (including generated `id` and timestamps).

### Mark as read

**Request:** `PATCH /messages/{id}/read` (no body).

**Response:** `200 OK` with the updated message.

## Project layout

```
cmd/api/              # Application entrypoint
internal/
  api/
    handlers/         # HTTP layer (Echo): bind JSON, call services, respond
    routes/           # Route registration
  config/             # Environment-based configuration
  domain/             # Models (GORM + JSON tags)
  repository/postgres # Data access (GORM only)
  service/            # Business logic and validation
```

Request flow: **routes → handlers → service → repository → domain**.

When you add a new persisted model, register it in `AutoMigrate` in `internal/repository/postgres/connection.go`.

## Development

```bash
# Format and static check
gofmt -w .
go vet ./...

# Tests
go test ./...
```

Do not commit `.env` or real credentials.
