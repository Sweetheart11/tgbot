build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

init_postgres:
	@docker-compose up --build -d
