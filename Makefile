.DEFAULT_GOAL := help
.PHONY: help proto build run dev clean deps install-tools docker-build docker-up docker-down docker-logs lint lint-fix

help:
	@echo "Available commands:"
	@echo "  make proto          - Generate protobuf code"
	@echo "  make build          - Build the application"
	@echo "  make run            - Build and run the application"
	@echo "  make dev            - Run in development mode (no build)"
	@echo "  make clean          - Remove build artifacts and data"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make install-tools  - Install protoc plugins"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start services with docker-compose"
	@echo "  make docker-down    - Stop services with docker-compose"
	@echo "  make docker-logs    - Show docker-compose logs"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make lint-fix       - Run golangci-lint with auto-fix"

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/imagestore/v1/imagestore.proto

build: proto
	go build -o bin/imagestore ./cmd/server

run: build
	./bin/imagestore

clean:
	rm -rf bin/
	rm -f imagestore.db
	find images -type f ! -name '.gitkeep' -delete

deps:
	go mod download
	go mod tidy

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)

dev:
	go run ./cmd/server

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...
