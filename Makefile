DB_CONTAINER=whoami-db-1
DB_NAME=whoami_db
DB_TEST_NAME=whoami_test
DB_USER=whoami_user
DB_PASS=secret
NETWORK_NAME=whoami_network
DB_URL=pgx5://${DB_USER}:${DB_PASS}@localhost:5432/${DB_NAME}?sslmode=disable
TEST_DB_URL=pgx5://${DB_USER}:${DB_PASS}@localhost:5432/${DB_TEST_NAME}?sslmode=disable
network:
	docker network create ${NETWORK_NAME}
postgres:
	docker run --name ${DB_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -d postgres:17-alpine
createdb:
	docker exec -it ${DB_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}
createdb_test:
	docker exec -it ${DB_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_TEST_NAME}
dropdb:
	docker exec -it ${DB_CONTAINER} dropdb ${DB_NAME}
migrateup:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up
migrateup_test:
	migrate -path internal/db/migrations -database "$(TEST_DB_URL)" -verbose up
migratedown:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down
migratedown_test:
	migrate -path internal/db/migrations -database "$(TEST_DB_URL)" -verbose down
new_migration:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run ./cmd/main.go

.PHONY: help
help:
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sed -e 's/\(.*\):.*##[ \t]*\(.*\)/  \1|\2/' | column -t -s '|' | sort

.DEFAULT_GOAL := help

.PHONY: network postgres createdb dropdb migrateup migratedown  new_migration sqlc test server
