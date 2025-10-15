FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git protobuf-dev make libwebp-dev gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/imagestore/v1/imagestore.proto

ENV CGO_ENABLED=1
RUN go build -o imagestore ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates libwebp

WORKDIR /root/

COPY --from=builder /app/imagestore .

RUN mkdir -p images/originals images/thumbnails

EXPOSE 50051 8087

CMD ["./imagestore"]
