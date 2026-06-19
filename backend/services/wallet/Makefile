##@ Format
.PHONY: Format
format:
	tool/scripts/format.sh

.PHONY: tidy
tidy:
	go mod tidy

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
pretty: tidy gen.db gen.mock format lint lint.struct

##@ Test
.PHONY: test.unit
test.unit:
	go clean -testcache
	go test -count=1 -failfast -v -race ./...

##@ Migration
.PHONY: migration
migration: ## Create database migration using atlas.
	tool/scripts/create-migration.sh wallet

.PHONY: migrate
migrate:
	atlas migrate apply --dir file://db/migrations --url "postgres://postgres:postgres@localhost:54322/postgres?sslmode=disable&search_path=wallet"

##@ Generator
.PHONY: gen.mock
gen.mock:
	mockery

.PHONY: gen.db
gen.db:
	sqlc generate
