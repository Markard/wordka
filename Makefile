ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
endif

# Exporting bin folder to the path for makefile
export PATH := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile


# ~~~ Builds ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
.PHONY: d-devenv-start
d-devenv-start: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker-compose up -d --build postgres

.PHONY: d-start
d-start: ## Start application silently
	docker-compose up -d --build

.PHONY: d-stop
d-stop: ## Stop application
	docker-compose down

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
install-deps: migrate air gotestsum tparse testfixtures ## Install Development Dependencies (localy).

deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(GOLANGCI) $(TESTFIXTURES) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-air: $(AIR) ## Starts AIR ( Continuous Development app).
	air

# ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
PG_DSN=postgres://${PG_USER}:${PG_PASS}@localhost:${PG_PORT}/${PG_DB}?sslmode=disable

.PHONY: migrate-up
migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
	@ migrate  -database $(PG_DSN) -path=migrations up

.PHONY: migrate-down
migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
 	migrate  -database $(PG_DSN) -path=migrations down $${N}

.PHONY: migrate-drop
migrate-drop: $(MIGRATE) ## Drop everything inside the database.
	migrate  -database $(PG_DSN) -path=migrations drop

.PHONY: migrate-create
migrate-create: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir migrations $${Name}

# ~~~ Test Fixtures ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
PATH_TO_FIXTURES=test/fixtures

.PHONY: fixtures-load
fixtures-load: $(TESTFIXTURES) ## Load a set of fixtures to database.
	testfixtures --dangerous-no-test-database-check -d postgres -c "$(PG_DSN)" -D $(PATH_TO_FIXTURES)

.PHONY: fixtures-dump
fixtures-dump: $(TESTFIXTURES) ## Dump a set of fixtures from the database.
	testfixtures --dangerous-no-test-database-check -d postgres -c "$(PG_DSN)" -D $(PATH_TO_FIXTURES)