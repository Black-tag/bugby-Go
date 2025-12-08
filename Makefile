# Makefile for Bugby-Go

# Variables
APP_NAME := bugby
BUILD_DIR := bin
MAIN_FILE := cmd/main.go
DOCKER_COMPOSE := docker compose

.PHONY: all build run test lint clean docker-up docker-down help

all: build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

## build: Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

## run: Run the application locally
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_FILE)

## test: Run tests with race detection
test:
	@echo "Running tests..."
	@go test -race -v ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@golangci-lint run --no-config --verbose



## clean: Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

## docker-up: Start Docker services
docker-up:
	@echo "Starting Docker services..."
	@$(DOCKER_COMPOSE) up -d --build

## docker-down: Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@$(DOCKER_COMPOSE) down
