.DEFAULT_GOAL := help

.PHONY: help proto build run dev clean test deps install-tools

help:
	@echo "Available commands:"
	@echo "  make proto          - Generate protobuf code"
	@echo "  make build          - Build the application"
	@echo "  make run            - Build and run the application"
	@echo "  make dev            - Run in development mode (no build)"
	@echo "  make clean          - Remove build artifacts and data"
	@echo "  make test           - Run tests"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make install-tools  - Install protoc plugins"

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/imagestore/v1/imagestore.proto

build: proto
	go build -o bin/imagestore .

run: build
	./bin/imagestore

clean:
	rm -rf bin/
	rm -f imagestore.db
	find images -type f ! -name '.gitkeep' -delete

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

dev:
	go run ./cmd/server
