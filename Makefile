.DEFAULT_GOAL := help
.PHONY: help proto dev install-tools prod-build prod-up prod-down lint lint-fix

help:
	@echo "Available commands:"
	@echo "  make proto          - Generate protobuf code"
	@echo "  make dev            - Run in development mode"
	@echo "  make install-tools  - Install protoc plugins"
	@echo "  make prod-build     - Build Docker image"
	@echo "  make prod-up        - Start services with docker-compose"
	@echo "  make prod-down      - Stop services with docker-compose"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make lint-fix       - Run golangci-lint with auto-fix"

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/imagestore/v1/imagestore.proto

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)

dev:
	go run ./cmd/server

prod-build:
	docker-compose build

prod-up:
	docker-compose up -d

prod-down:
	docker-compose down

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...
