package internal

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type Storage struct {
	baseDir       string
	thumbnailSize uint
}

func NewStorage(baseDir string) *Storage {
	return &Storage{
		baseDir:       baseDir,
		thumbnailSize: 200,
	}
}

func (s *Storage) SaveImage(userID, filename string, reader io.Reader) (imageID, originalPath, thumbnailPath string, width, height int32, size int64, err error) {
	imageID = uuid.New().String()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("failed to read image: %w", err)
	}
	size = int64(len(data))

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width = int32(bounds.Dx())
	height = int32(bounds.Dy())

	originalDir := filepath.Join(s.baseDir, "originals", userID)
	if err := os.MkdirAll(originalDir, 0o755); err != nil {
		return "", "", "", 0, 0, 0, err
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = "." + format
	}
	originalFullPath := filepath.Join(originalDir, imageID+ext)

	if err := os.WriteFile(originalFullPath, data, 0o644); err != nil {
		return "", "", "", 0, 0, 0, err
	}

	thumbnailDir := filepath.Join(s.baseDir, "thumbnails", userID)
	if err := os.MkdirAll(thumbnailDir, 0o755); err != nil {
		os.Remove(originalFullPath)
		return "", "", "", 0, 0, 0, err
	}

	thumbnail := resize.Thumbnail(s.thumbnailSize, s.thumbnailSize, img, resize.Lanczos3)
	thumbnailFullPath := filepath.Join(thumbnailDir, imageID+"_thumb.webp")

	thumbFile, err := os.Create(thumbnailFullPath)
	if err != nil {
		os.Remove(originalFullPath)
		return "", "", "", 0, 0, 0, err
	}
	defer thumbFile.Close()

	if err := webp.Encode(thumbFile, thumbnail, &webp.Options{Lossless: false, Quality: 80}); err != nil {
		os.Remove(originalFullPath)
		os.Remove(thumbnailFullPath)
		return "", "", "", 0, 0, 0, err
	}

	originalPath = filepath.Join("originals", userID, imageID+ext)
	thumbnailPath = filepath.Join("thumbnails", userID, imageID+"_thumb.webp")

	return imageID, originalPath, thumbnailPath, width, height, size, nil
}

func (s *Storage) DeleteImage(originalPath, thumbnailPath string) error {
	os.Remove(filepath.Join(s.baseDir, originalPath))
	os.Remove(filepath.Join(s.baseDir, thumbnailPath))
	return nil
}

func (s *Storage) GetImagePath(imagePath string) string {
	return filepath.Join(s.baseDir, imagePath)
}
