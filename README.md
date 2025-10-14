# Image Storage API

Go-based image storage service with gRPC and HTTP APIs. Supports streaming uploads, automatic thumbnail generation, and user-based organization.

## Features

- gRPC streaming uploads + HTTP image serving
- Automatic WebP thumbnail generation (200x200px)
- SQLite metadata storage + filesystem image storage
- User-based organization with pagination
- Domain-Driven Design architecture

## Quick Start

```bash
# Install protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Setup and run
./setup.sh
./bin/imagestore
```

**Servers:**
- gRPC: `localhost:50051`
- HTTP: `localhost:8080`

**Or use the automated script:**
```bash
chmod +x setup.sh && ./setup.sh && ./bin/imagestore
```

## API

### gRPC (port 50051)
- `UploadImage(stream)` - Upload image with streaming
- `GetImageMetadata(image_id)` - Get image info
- `ListImages(user_id, page_size, page_token)` - List user images
- `DeleteImage(image_id, user_id)` - Delete image
- `GetImageURL(image_id, thumbnail)` - Get HTTP URL

### HTTP (port 8080)
- `GET /images/{id}` - Serve original image
- `GET /images/{id}?thumbnail=true` - Serve thumbnail
- `GET /health` - Health check

## Example Usage

**Build and use the example client:**
```bash
go build -o bin/client example_client.go

# Upload
./bin/client -action=upload -user=user123 -file=photo.jpg

# List
./bin/client -action=list -user=user123

# Get metadata
./bin/client -action=get -id=<image-id>

# Delete
./bin/client -action=delete -user=user123 -id=<image-id>
```

**Go client code:**
```go
conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
client := pb.NewImageServiceClient(conn)

stream, _ := client.UploadImage(context.Background())

// Send metadata
stream.Send(&pb.UploadImageRequest{
    Data: &pb.UploadImageRequest_Metadata{
        Metadata: &pb.ImageMetadataInput{
            UserId:      "user123",
            Filename:    "photo.jpg",
            ContentType: "image/jpeg",
        },
    },
})

// Send chunks
buf := make([]byte, 64*1024)
for {
    n, err := file.Read(buf)
    if err == io.EOF { break }
    stream.Send(&pb.UploadImageRequest{
        Data: &pb.UploadImageRequest_Chunk{Chunk: buf[:n]},
    })
}

resp, _ := stream.CloseAndRecv()
fmt.Printf("URL: %s\n", resp.Url)
```

**Display in browser:**
```html
<img src="http://localhost:8080/images/{image_id}">
<img src="http://localhost:8080/images/{image_id}?thumbnail=true">
```

## Architecture

**Simple Layered Structure:**

```
imagestore/
├── main.go          # Server setup & wiring
├── handlers.go      # All gRPC + HTTP handlers
├── storage.go       # Image storage operations
├── db.go           # Database operations
├── types.go        # Shared structs
├── proto/          # Protobuf definitions
├── images/         # Storage
│   ├── originals/  # {user_id}/{image_id}.ext
│   └── thumbnails/ # {user_id}/{image_id}_thumb.webp
├── go.mod
└── Makefile
```

**Just 5 Go files - simple and easy to understand!**

## Configuration

Edit `cmd/server/main.go`:
```go
const (
    grpcPort  = ":50051"
    httpPort  = ":8080"
    baseURL   = "http://localhost:8080"
    dbPath    = "./imagestore.db"
    imagesDir = "./images"
)
```

## Integration with Your App

1. Copy `proto/imagestore/v1/imagestore.proto` to your project
2. Generate client: `protoc --go_out=. --go-grpc_out=. imagestore.proto`
3. Update module path: Replace `github.com/yourusername/imagestore` in `go.mod` and all imports
4. Connect and use (see example above)

## Docker

```bash
docker-compose up -d
```

## Dependencies

- `google.golang.org/grpc` - gRPC framework
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/google/uuid` - UUID generation
- `github.com/chai2010/webp` - WebP encoding
- `github.com/nfnt/resize` - Image resizing

## Notes

- Images stored in filesystem under `images/` directory
- Thumbnails automatically generated as 200x200 WebP
- SQLite handles metadata (dimensions, upload time, etc)
- HTTP URLs can be embedded directly in HTML for browser display
- Authorization: basic user_id check (enhance for production)

## License

MIT
