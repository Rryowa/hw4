# Variables
BINARY_NAME=cli
BIN_DIR=bin
CMD_DIR=cmd
EXPLAIN_DIR=explain
DB_STRING="postgres://avrigne:8679@localhost/cli?sslmode=disable"
MIGRATIONS_DIR="./migrations"

up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down

build:
	@echo "Building the CLI application..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@echo "Build completed. Binary is located at $(BIN_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running the CLI application..."
	@$(BIN_DIR)/$(BINARY_NAME)

.PHONY: up down build run