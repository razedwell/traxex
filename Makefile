PROTO_DIR = proto/definitions
PROTO_GEN = proto/gen/go

PROTO_FILES = $(shell find $(PROTO_DIR) -name '*.proto')

DB_DSN ?=postgres://traxex:drowranger@localhost:5432/traxex_dev?sslmode=disable
MIGRATIONS_DIR := services/auth/migrations

.PHONY: proto

generate: gen-proto gen-mocks gen-openapi

gen-proto: 
	@echo "Generating Go code from proto files..."
	@mkdir -p $(PROTO_GEN)
	protoc \
		-I=$(PROTO_DIR) \
		--go_out=$(PROTO_GEN) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

gen-mocks:
	@echo "Generating mocks..."
	mockery --config .mockery.yaml

gen-openapi:
	@echo "Generating OpenAPI code..."
	oapi-codegen -config services/auth/oapi-codegen.yaml services/auth/api/openapi3/openapi.yaml

test:
	@echo "Running tests..."
	go test -v ./...

build:
	@echo "Building the binary artifact..."
	go build -o bin/auth-service ./services/auth/cmd/server

docker-up:
	@echo "Starting services with Docker Compose..."
	docker compose \
		--env-file deployments/.env \
		-f deployments/docker-compose.yaml up -d --build

docker-down:
	@echo "Stopping services with Docker Compose..."
	docker compose \
		--env-file deployments/.env \
		-f deployments/docker-compose.yaml down -v

docker-build:
	@echo "Building Docker images..."
	docker build -t auth-service:latest -f services/auth/Dockerfile .

migrate-create:
	@echo "Creating a new migration file: $(name) ..."
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	@echo "Migrating up..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	@echo "Migrating down..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

clean:
	@echo "Deleting generated files..."
	@rm -rf $(PROTO_GEN)/*
	@rm -rf bin/*

deps:
	@echo "Installing dependencies..."
	@echo "Installing protocol buffer compiler and Go plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Dependencies installed successfully."

lint:
	@echo "Linting the code..."
	go vet ./...