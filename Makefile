.PHONY: proto build run clean test

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
