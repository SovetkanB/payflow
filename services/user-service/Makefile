run:
	go run ./cmd/main.go
build:  # go build -o bin/user-service ./cmd/main.go
test:   # go test ./...

DB_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable

MIGRATIONS_DIR=./migrations

migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -all

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version