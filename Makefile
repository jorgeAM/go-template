include .env

# Go
generate:
	go generate ./...

unit-tests:
	go test ./... -cover -v -coverprofile=./coverage.out

integration-tests:
	@echo "integration-tests: to be implemented soon"

show-cover:
	go tool cover -html=./coverage.out

tidy:
	go mod tidy
	go mod vendor

run:
	go run ./cmd/app | jq '.'

# Migrations
# `make migrate` applies every module's pending up migrations (cmd/migrate walks
# the same module list cmd/app boots from; each module keeps its own
# schema_migrations_<name> table). new_migration and migration_down need
# MODULE=<name> and the golang-migrate CLI.
MIGRATIONS_DIR = internal/$(MODULE)/infrastructure/persistence/migrations
MIGRATIONS_TABLE = schema_migrations_$(MODULE)
postgres_url = "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable&x-migrations-table=${MIGRATIONS_TABLE}"

migrate:
	go run ./cmd/migrate

new_migration: check-module
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATION_NAME)

migration_down: check-module
	migrate -path $(MIGRATIONS_DIR) -database $(postgres_url) -verbose down

check-module:
ifndef MODULE
	$(error MODULE is required, e.g. make new_migration MODULE=identity MIGRATION_NAME=add_name)
endif
