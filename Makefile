start:
	docker-compose up --build

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

.PHONY: start build lint test test-coverage install create-migration down-migration
