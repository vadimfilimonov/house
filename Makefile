start:
	go run cmd/main.go

start-db:
	redis-server

build:
	go build -o houseBuild cmd/main.go
	chmod +x houseBuild

lint:
	go vet ./...

test:
	go test ./...

test-coverage:
	go test ./... -cover

install:
	go mod tidy

# Название директории, куда попадут файлы миграции
MIGRATE_DIR=schema

create-migration:
	@read -p "Введите название таблицы: " NAME; \
	migrate create -ext sql -dir $(MIGRATE_DIR) $$NAME

.PHONY: start build lint test test-coverage install create-migration
