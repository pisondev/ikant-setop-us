API_DIR := apps/api
MIGRATIONS_DIR := $(API_DIR)/migrations

DB_HOST ?= localhost
DB_PORT ?= 5438
DB_USER ?= ikant_user
DB_PASSWORD ?= ikant_pass
DB_NAME ?= ikant_setop_us_db
DB_SSLMODE ?= disable
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

MIGRATE := go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

.PHONY: migrate-version-windows migrate-up-windows migrate-down-windows migrate-down-all-windows migrate-reset-windows migrate-force-windows migrate-create-windows

migrate-version-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-up-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-down-all-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

migrate-reset-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" drop -f
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-force-windows:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)

migrate-create-windows:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
