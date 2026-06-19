include Makefile.help.mk

##@ Format
.PHONY: Format
format: ## Format codes using goimports and gofmt.
	tool/scripts/format.sh

.PHONY: tidy
tidy: ## Run `go mod tidy`.
	go mod tidy

.PHONY: lint.sql
lint.sql: ## Lint sql files.
	sqlfluff fix --show-lint-violations -d postgres -e LT05,AM04,AM09,PG01

.PHONY: lint.struct
lint.struct: ## Run fieldalignment.
	tool/scripts/lint-struct.sh

.PHONY: lint.cleancache
lint.cleancache: ## Clean golangci-lint cache.
	golangci-lint cache clean

.PHONY: lint
lint: ## Lint proto files using buf and golang files using golangci-lint.
lint: lint.cleancache
	tool/scripts/lint.sh

.PHONY: pretty
pretty: ## Prettify golang and proto files. Basically, it runs tidy, format, and lint command.
pretty: tidy gen.db gen.mock format lint.sql lint lint.struct

.PHONY: check.import
check.import: ## Check if import blocks are separated accordingly.
	tool/scripts/check-import.sh

##@ Test
.PHONY: test.unit
test.unit: ## Run unit test.
	go clean -testcache
	go test -count=1 -failfast -v -race ./...

.PHONY: test.cover
test.cover: ## Run unit test and generates coverage file.
	go clean -testcache
	go test -count=1 -failfast -v -race -coverprofile=coverage.out $$(go list ./... | grep -v "/test")
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

##@ Migration
.PHONY: migration
migration: ## Create database migration using atlas.
	tool/scripts/create-migration.sh $(name)

.PHONY: migrate
migrate: ## Apply migration.
	atlas migrate apply --env local

##@ Generator
.PHONY: gen.mock
gen.mock: ## Generate mocks.
	mockery

.PHONY: gen.db
gen.db: ## Generate DB implementation using sqlc.
	sqlc generate
