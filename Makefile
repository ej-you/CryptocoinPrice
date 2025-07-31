go_exec="./cmd/app/main.go"
server_runner_path="./internal/app/server/server.go"
go_migrator_path="./cmd/migrator/main.go"

# title of migration
title = "migration"
version = 1

# --- #
# APP #
# --- #

dev:
	go run $(go_exec)

lint:
	golangci-lint run -c ./.golangci.yml ./...

# ------- #
# SWAGGER #
# ------- #

swagger-update:
	@swag init -g $(server_runner_path)

swagger-fmt:
	@swag fmt -g $(server_runner_path)

swagger: swagger-fmt swagger-update

# ---------- #
# MIGRATIONS #
# ---------- #

# use "title" var for name the migration
.PHONY: migrations
migrations:
	migrate create -digits 2 -dir migrations -ext sql -seq "$(title)"

# use "version" var for specify the version of the migration to force
migrate-force:
	@go run $(go_migrator_path) force $(version)

migrate-status:
	@go run $(go_migrator_path) status

migrate-up:
	@go run $(go_migrator_path) up -n 1

migrate-down:
	@go run $(go_migrator_path) down -n 1
