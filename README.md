# House Service

Backend service for the Avito Backend Bootcamp house and flat moderation task.

## Documentation

- [Task description](doc/assignment.md)
- [System documentation](doc/docs.md)
- [OpenAPI specification](api.yaml)
- [Project structure diagram](doc/project-structure.svg)

## Local Run

Start the complete development environment in one command:

```bash
make run
```

This starts the Go service, PostgreSQL and Redis. The service is available at
`http://localhost:8080`, PostgreSQL at `localhost:5432`, and Redis at
`localhost:6379`. Database migrations run automatically when the service starts.

Useful commands:

```bash
make dev-logs
make db-shell
make dev-down
```

The PostgreSQL data is stored in the Docker volume `house_db_data` and survives
container restarts. To remove it and recreate the database from scratch:

```bash
docker compose down -v
```

To run the Go service directly on the host, start only PostgreSQL and Redis and
set `DB_HOST=localhost` and `REDIS_HOST=localhost` in `.env`:

```bash
docker compose up -d db redis
go run ./cmd
```

## Production Docker Image

The production image does not include the source code, hot reload tools, or a
local database. It expects PostgreSQL and Redis to be provided by the hosting
platform or managed services.

Create the production environment file and fill in real credentials and hostnames:

```bash
cp .env.prod.example .env.prod
make deploy-prod
```

The production image can also be built separately:

```bash
make prod-build
```

`make prod-run` is an alias for `make deploy-prod`.

Stop the production container with:

```bash
make prod-down
```

Never commit `.env.prod` or put production secrets into the repository.

## Checks

```bash
go test ./...
go vet ./...
go build -o /tmp/house-service ./cmd
```
