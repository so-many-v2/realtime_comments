-include .env
export

MIGRATIONS_DIR := migrations
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=disable

make_migrations_template:
	@test -n "$(name)" || (echo "usage: make make_migrations_template name=<migration_name>" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate_down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

run_compose:
	docker-compose up --build -d

reup:
	docker-compose down
	docker-compose up --build -d

reset_database:
	docker-compose down -v
	docker-compose up --build -d


