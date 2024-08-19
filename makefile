.PHONY: up

MIGRATE_CMD=migrate
MIGRATE_DIR=./migrations
DB_DSN="mysql://multifinance-credit-user:multifinance-credit-pw@tcp(localhost:3306)/multifinance-credit-db?parseTime=true"
DATE=$(shell date +%Y%m%d_%H%M%S)

generate: api/api.yml
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
	docker compose -f deployment/docker-compose.yml up -d

stop-infra:
	docker compose -f deployment/docker-compose.yml stop