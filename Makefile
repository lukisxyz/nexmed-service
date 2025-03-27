# Variables
APP_NAME=nexmedis
DB_URL ?= postgres://fahmi:wap12345@localhost:5432/nexmedis?sslmode=disable

# Go commands
run:
	go run main.go

build:
	go build -o bin/$(APP_NAME) main.go

test:
	go test -v ./...

# Database migrations
migrate-up:
	@if [ "$(db)" ]; then \
		migrate -database "$(db)" -path migrations up; \
	else \
		migrate -database "$(DB_URL)" -path migrations up; \
	fi

migrate-down:
	@if [ "$(db)" ]; then \
		migrate -database "$(db)" -path migrations down; \
	else \
		migrate -database "$(DB_URL)" -path migrations down; \
	fi

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Install tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install swagger tools
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger documentation
swagger-gen:
	swag init -g main.go -o docs

# Generate swagger with specific configuration
swagger-gen-config:
	swag init \
		-g main.go \
		-o docs \
		--parseDependency \
		--parseInternal \
		--parseDepth 1 \
		--instanceName swagger

# Serve swagger UI locally
swagger-serve:
	go run main.go -docs

# Combined swagger command
swagger: swagger-install swagger-gen

# Clean swagger files
swagger-clean:
	rm -rf docs/swagger.*

.PHONY: run build test migrate-up migrate-down migrate-create install-tools docker-build docker-run swagger swagger-install swagger-gen swagger-serve
