-include .env
export

BINARY     := bin/app
CMD        := ./cmd/api
MIGRATIONS := ./migrations

DB_URL     ?= $(DATABASE_URL)
MIGRATE    := migrate -path $(MIGRATIONS) -database "$(DB_URL)"

.PHONY: run build migrate-up migrate-down migrate-drop migrate-create

run:
	go run $(CMD)/main.go

build:
	go build -o $(BINARY) $(CMD)/main.go

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-drop:
	$(MIGRATE) drop -f

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS) -seq $$name
