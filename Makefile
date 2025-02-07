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

.PHONY: start build lint test test-coverage install
