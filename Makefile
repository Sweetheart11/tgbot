include .env
export 

build:	
	@go build -o bin/gobank ./cmd/main.go

run: build
	@./bin/gobank

init_postgres:
	@docker-compose up --build -d

up_migrations: 
	@cd storage/migrations ; \
	goose postgres "host=${POSTGRESDB_HOST} user=${POSTGRESDB_USER} database=${POSTGRESDB_NAME} password=${POSTGRESDB_PASSWORD} sslmode=${POSTGRESDB_SSLMODE}" up

down_migrations: 
	@cd storage/migrations ; \
	goose postgres "host=${POSTGRESDB_HOST} user=${POSTGRESDB_USER} database=${POSTGRESDB_NAME} password=${POSTGRESDB_PASSWORD} sslmode=${POSTGRESDB_SSLMODE}" down
