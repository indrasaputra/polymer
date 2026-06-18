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


##@ Test
.PHONY: test.unit
test.unit:
	go clean -testcache
	go test -count=1 -failfast -v -race ./...

##@ Migration
.PHONY: migration
migration: ## Create database migration using atlas.
	tool/scripts/create-migration.sh wallet

##@ Generator
.PHONY: gen.mock
gen.mock:
	mockery
