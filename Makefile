.PHONY: run test lint build migrations/new migrations/up migrations/down
run:
	go run ./cmd/api
test:
	go test -race ./...
lint:
	go vet ./...
build:
	go build -o bin/api ./cmd/api
migrations/new:
	@test -n "$(NAME)" || (echo "NAME is required"; exit 1)
	migrate create -seq -digits 6 -ext sql -dir ./migrations "$(NAME)"
migrations/up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up
migrations/down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1
