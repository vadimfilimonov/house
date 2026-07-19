# House Service

Backend service for the Avito Backend Bootcamp house and flat moderation task.

## Documentation

- [Task description](doc/assignment.md)
- [OpenAPI specification](api.yaml)
- [Project structure diagram](doc/project-structure.svg)

## Local Run

Copy environment variables and start dependencies:

```bash
cp .env.example .env
docker compose up -d
```

Run the service:

```bash
go run ./cmd
```

## Checks

```bash
go test ./...
go vet ./...
go build ./cmd
```
