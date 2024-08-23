.PHONY: up

MIGRATE_CMD=migrate
MIGRATE_DIR=./migrations
DB_DSN="mysql://multifinance-credit-user:multifinance-credit-pw@tcp(localhost:3306)/multifinance-credit-db?parseTime=true"
DATE=$(shell date +%Y%m%d_%H%M%S)

install:
	go mod download
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.14
	go install go.uber.org/mock/mockgen@v0.3.0
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2
	npm install -g @redocly/cli@1.9.1

# Generates mocks for interfaces
INTERFACES_GO_FILES := $(shell find internal -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

generate: api/api.yml generate_mocks
	mkdir -p generated/api
	oapi-codegen --package api -generate types $< > generated/api/api-types.gen.go
create:
	@$(MIGRATE_CMD) create -ext sql -dir $(MIGRATE_DIR) $(NAME)

up:
	@$(MIGRATE_CMD) -source file://$(MIGRATE_DIR) -database=$(DB_DSN) up

force:
	@$(MIGRATE_CMD) -source file://$(MIGRATE_DIR) -database=$(DB_DSN) force 20240819170341

reset:
	@$(MIGRATE_CMD) reset -dir $(MIGRATE_DIR)

refresh: reset up

down:
	@$(MIGRATE_CMD) -source file://$(MIGRATE_DIR) -database=$(DB_DSN) down

status:
	@$(MIGRATE_CMD) status -dir $(MIGRATE_DIR)

starting-infra:
	docker compose -f docker-compose.yml up -d

stop-infra:
	docker compose -f docker-compose.yml stop

test:
	go test -tags test -short -failfast -coverprofile coverage.out ./...

init: starting-infra \
	install \
	generate \
	up \
	run

run:
	air -c air.toml

clean:
	rm -rf dist/* generated build vendor
	find . -name "*.mock.gen.go" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "wire_gen.go" -type f -delete
	find . -name "*.mock.gen.go" -type f -delete

preview_open_api:
	redocly preview-docs api/api.yml