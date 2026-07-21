start:
	docker compose up --build

run:
	docker compose up --build -d

dev:
	$(MAKE) run

prod-build:
	docker build -f Dockerfile -t $${IMAGE_NAME:-house:prod} .

deploy-prod:
	@test -f .env.prod || (echo 'Create .env.prod from .env.prod.example first'; exit 1)
	docker compose --env-file .env.prod -f docker-compose.prod.yml up --build -d

prod-run:
	$(MAKE) deploy-prod

prod-down:
	docker compose --env-file .env.prod -f docker-compose.prod.yml down

dev-down:
	docker compose down

dev-logs:
	docker compose logs -f app

db-shell:
	docker compose exec db sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

start-db:
	redis-server

build:
	go build -o houseBuild cmd/main.go
	chmod +x houseBuild

lint:
	docker run --rm -v $(PWD):/app -w /app house-app go vet ./...

test:
	docker run --rm -v $(PWD):/app -w /app house-app go test ./...

test-coverage:
	go test ./... -cover

install:
	go mod tidy

# Название директории, куда попадут файлы миграции
MIGRATE_DIR=schema

create-migration:
	@read -p "Введите название таблицы: " NAME; \
	migrate create -ext sql -dir $(MIGRATE_DIR) $$NAME

down-migration:
	@read -p "Введите строку подключения к БД: " DB_PATH; \
	migrate -path ${MIGRATE_DIR} -database "$$DB_PATH" down --all

.PHONY: start run dev deploy-prod prod-run prod-down dev-down dev-logs db-shell build lint test test-coverage install create-migration down-migration
