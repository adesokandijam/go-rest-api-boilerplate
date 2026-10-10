.PHONY: run test lint build
run:
	go run ./cmd/api
test:
	go test -race ./...
lint:
	go vet ./...
build:
	go build -o bin/api ./cmd/api
