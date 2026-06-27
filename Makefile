-include .env
export

MIGRATIONS_DIR := migrations
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=disable

# Postgres Cmd
make_migrations_template:
	@test -n "$(name)" || (echo "usage: make make_migrations_template name=<migration_name>" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate_down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

clear_db:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) < $(MIGRATIONS_DIR)/sql/clear_db.sql

setup_db:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) < $(MIGRATIONS_DIR)/sql/setup_db.sql

# Compose cmd
run_compose:
	docker-compose up --build -d

reup:
	docker-compose down
	docker-compose up --build -d


# K6 cmd
load_post:
	k6 run tests/load/post_service.js

load_comment:
	k6 run tests/load/comment_service.js

stress_comment:
	k6 run tests/load/comment_stress.js


