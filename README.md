# Go REST API Boilerplate

A small, dependency-free starting point for production Go HTTP APIs.

## Included

- Go 1.25 `net/http` routing
- JSON responses and versioned routes
- Structured request logging, request IDs, panic recovery, and timeouts
- Graceful shutdown and environment-based configuration
- Unit tests, race detection, vet, GitHub Actions, Docker, and MIT license

## Run

```sh
cp .env.example .env
go run ./cmd/api
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/hello/world
```

## Commands

```sh
make run
make test
make lint
make build
```

Rename the module in `go.mod`, then add domain packages under `internal/` as the application grows.
