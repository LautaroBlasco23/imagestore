# Image Storage Service

I created this project to learn Go, gRPC, and image storage. Feel free to use it however you like.

# What is it?

A lightweight image storage service built with Go. Upload images via gRPC, serve them over HTTP, and get automatic WebP thumbnails.

## Service Summary

This service handles the entire lifecycle of user-uploaded images:
- Accepts image uploads through a streaming gRPC interface
- Validates and processes them
- Generates optimized thumbnails
- Tracks metadata in SQLite
- Serves both originals and thumbnails through HTTP API

## Getting Started

First, install the protobuf tools:
```bash
make install-tools
```

Generate the gRPC code:
```bash
make proto
```

### Development

Run locally for development:
```bash
make dev
```

Your service will be ready at:
- gRPC on `localhost:50051`
- HTTP on `localhost:8087`

Check it's working:
```bash
curl http://localhost:8087/health
```

### Production

Build and run with Docker:
```bash
make prod-build
make prod-up
```

Stop everything:
```bash
make prod-down
```

Images persist in the `imagestore_data` Docker volume, so they survive restarts.

## Storage Structure

Images are organized in a directory hierarchy:
- Originals: `baseDir/originals/userID/{imageID}.ext`
- Thumbnails: `baseDir/thumbnails/userID/{imageID}_thumb.webp`

Each image gets a unique UUID-based identifier. Original images preserve their format, thumbnails are always WebP.

## Thumbnail Generation

- **Size**: 200x200px
- **Algorithm**: Lanczos3 resampling (high-quality downscaling)
- **Format**: WebP with 80% quality
- **Cleanup**: If any step fails, previously created files are removed

## Security Features

### Path Traversal Protection

- Ensures all file paths stay within `baseDir` boundaries
- Prevents attacks like `../../etc/passwd` by cleaning paths and checking prefixes
- Validates separator boundaries to prevent edge cases

### Safe Extension Handling

- Limits extension length to 5 characters
- Validates extension format (must start with `.`)
- Falls back to image format if filename extension is suspicious

### File Permissions

- Directories: `0o750` (rwxr-x---)
- Files: `0o600` (rw-------)

## Performance Considerations

The service is solid for moderate traffic but could benefit from streaming for large images and background processing for thumbnail generation at scale.

- Lanczos3 resampling is high quality but computationally expensive
- WebP encoding with 80% quality balances size/quality
- Synchronous operations block until complete

## License

MIT License - feel free to use this code for any purpose.
