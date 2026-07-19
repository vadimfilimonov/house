# House Service

Backend service for the Avito Backend Bootcamp house and flat moderation task.

## Documentation

- [Task description](doc/assignment.md)
- [System documentation](doc/docs.md)
- [OpenAPI specification](api.yaml)
- [Project structure diagram](doc/project-structure.svg)

## Local Run

Run everything in Docker Compose:

```bash
cp .env.example .env
docker compose up -d
```

Or run only PostgreSQL and Redis in Docker Compose, then start the service locally.
For local `go run`, set `DB_HOST=localhost` and `REDIS_HOST=localhost` in `.env`:

```bash
docker compose up -d db redis
go run ./cmd
```

## Checks

```bash
go test ./...
go vet ./...
go build -o /tmp/house-service ./cmd
```
