include .env

# Go
generate:
	go generate ./...

test:
	go test ./... -cover -v -coverprofile=./coverage.out

show-cover:
	go tool cover -html=./coverage.out

tidy:
	go mod tidy
	go mod vendor

run:
	go run cmd/app/*.go

# Migrations
postgres_url = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

new_migration:
	migrate create -ext sql -dir database/migration/ -seq $(MIGRATION_NAME)

migration_up:
	migrate -path database/migration/ -database $(postgres_url) -verbose up

migration_down:
	migrate -path database/migration/ -database $(postgres_url) -verbose down
