build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

init_postgres:
	@docker run --name tgbotdb -e POSTGRES_PASSWORD=G20rcpyxCD8 -d postgres
