# Variables
APP_NAME=nexmedis
DB_URL ?= postgres://fahmi:wap12345@localhost:5432/nexmedis?sslmode=disable
BUILD_DIR=build

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Air configuration
AIR_CONFIG=.air.toml

# Development server
dev:
	@if ! which air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air -c $(AIR_CONFIG)

# Build the application
build:
	@echo "Building..."
	@go build -o $(GOBIN)/$(APP_NAME) ./cmd/main.go

# Clean build directory
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
# Go commands
run:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) main.go

test:
	go test -v ./...

# Database migrations
migrate-up:
	@if [ "$(db)" ]; then \
		migrate -database "$(db)" -path db/migrations up; \
	else \
		migrate -database "$(DB_URL)" -path db/migrations up; \
	fi

migrate-down:
	@if [ "$(db)" ]; then \
		migrate -database "$(db)" -path db/migrations down; \
	else \
		migrate -database "$(DB_URL)" -path db/migrations down; \
	fi

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

# Install tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install swagger tools
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger documentation
swagger-gen:
	swag init -g cmd/main.go -o docs

# Generate swagger with specific configuration
swagger-gen-config:
	swag init \
		-g cmd/main.go \
		-o docs \
		--parseDependency \
		--parseInternal \
		--parseDepth 1 \
		--instanceName swagger

# Serve swagger UI locally
swagger-serve:
	go run cmd/main.go -docs

# Combined swagger command
swagger: swagger-install swagger-gen

# Clean swagger files
swagger-clean:
	rm -rf docs/swagger.*

.PHONY: run build test migrate-up migrate-down migrate-create install-tools docker-build docker-run swagger swagger-install swagger-gen swagger-serve
