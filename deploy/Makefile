include .env

LOCAL_BIN := $(CURDIR)/bin
MIGRATION_DIR := ./migrations
MIGRATION_DSN := "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

.PHONY: install-deps local-migration-up local-migration-down local-migration-status

install-deps:
	@mkdir -p "$(LOCAL_BIN)"
	@GOBIN="$(LOCAL_BIN)" go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	@"$(LOCAL_BIN)/goose" -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) status -v

local-migration-up:
	@"$(LOCAL_BIN)/goose" -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) up -v

local-migration-down:
	@"$(LOCAL_BIN)/goose" -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) down -v
