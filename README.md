# Image Storage API

Go-based image storage service with gRPC and HTTP APIs for uploading, serving, and managing images.

## Features

- gRPC streaming uploads
- HTTP image serving
- Automatic WebP thumbnail generation (200x200px)
- SQLite metadata storage
- User-based organization

## Quick Start

**Install tools:**
```bash
make install-tools
```

**Generate protobuf code:**
```bash
make proto
```

**Development:**
```bash
make dev
```

Servers will run on:
- gRPC: `localhost:50051`
- HTTP: `localhost:8087`

**Production:**
```bash
make prod-build
make prod-up
make prod-down  # to stop
```

## Available Commands

```bash
make help          # Show all available commands
make proto         # Generate protobuf code
make dev           # Run in development mode
make install-tools # Install protoc plugins
make prod-build    # Build Docker image
make prod-up       # Start services with docker-compose
make prod-down     # Stop services with docker-compose
make lint          # Run golangci-lint
make lint-fix      # Run golangci-lint with auto-fix
```

## API

### gRPC (port 50051)
- `UploadImage(stream)` - Upload image with streaming
- `GetImageMetadata(image_id)` - Get image info
- `ListImages(user_id, page_size, page_token)` - List user images
- `DeleteImage(image_id, user_id)` - Delete image
- `GetImageURL(image_id, thumbnail)` - Get HTTP URL

### HTTP (port 8087)
- `GET /images/{id}` - Serve original image
- `GET /images/{id}?thumbnail=true` - Serve thumbnail
- `GET /health` - Health check

## Storage

Images are stored in the `images/` directory:
- Originals: `images/originals/{user_id}/{image_id}.ext`
- Thumbnails: `images/thumbnails/{user_id}/{image_id}_thumb.webp`

## Docker

Services and volumes configured in `docker-compose.yml`. Images persist in the `imagestore_data` volume.
