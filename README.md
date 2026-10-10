# Go REST API Boilerplate

A production-minded Go REST API starter organised using the same application pattern as the Gym Tracker API.

## Included

- Flat `cmd/api` package for configuration, routes, handlers, middleware, helpers, and errors
- `internal/data` PostgreSQL model layer and `internal/validator`
- Complete example CRUD resource with optimistic versioning
- Versioned SQL migrations
- Structured logging, panic recovery, CORS, request limits, timeouts, and graceful shutdown
- Docker, Make targets, race-tested CI, vet, and an MIT license

## Project layout

```text
cmd/api/            server entry point, routes, handlers, middleware
internal/data/      database models and queries
internal/validator/ reusable validation helpers
migrations/         versioned PostgreSQL migrations
```

## Quick start

```sh
cp .env.example .env
createdb app
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/app?sslmode=disable'
make migrations/up
make run
```

Try the API:

```sh
curl http://localhost:8080/v1/healthcheck
curl -X POST http://localhost:8080/v1/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"Example","description":"Replace this resource with your own"}'
```

## Development

```sh
make test
make lint
make build
make migrations/new NAME=create_users_table
```

Rename the module in `go.mod`, replace the example `Item` resource, and add the authentication or other infrastructure required by your application.
