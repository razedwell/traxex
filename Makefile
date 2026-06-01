PROTO_DIR = proto/definitions
PROTO_GEN = proto/gen/go

PROTO_FILES = $(shell find $(PROTO_DIR) -name '*.proto')

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
	go build -o bin/auth-service ./services/auth/cmd

docker-compose:
	@echo "Starting services with Docker Compose..."
	docker compose up -d --build -f deploy/docker-compose.yaml

docker-build:
	@echo "Building Docker images..."
	docker build -t auth-service:latest -f services/auth/Dockerfile .

migrate-up:
	@echo "TODO migrate-up"

migrate-down:
	@echo "TODO migrate-down"

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
	@echo "Dependencies installed successfully."
