start:
	go run cmd/main.go

build:
	go build -o houseBuild cmd/main.go

lint:
	go vet ./...

test:
	go test ./...

test-coverage:
	go test ./... -cover

install:
	go mod tidy

.PHONY: start build lint test test-coverage install
