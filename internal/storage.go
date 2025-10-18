package internal

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

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
	dx, dy := bounds.Dx(), bounds.Dy()
	if dx < 0 || dy < 0 || dx > math.MaxInt32 || dy > math.MaxInt32 {
		return "", "", "", 0, 0, 0, fmt.Errorf("invalid image dimensions")
	}
	width = int32(dx)
	height = int32(dy)

	originalDir := filepath.Join(s.baseDir, "originals", userID)
	if err := os.MkdirAll(originalDir, 0o750); err != nil {
		return "", "", "", 0, 0, 0, err
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = "." + format
	}
	originalFullPath := filepath.Join(originalDir, imageID+ext)
	if err := os.WriteFile(originalFullPath, data, 0o600); err != nil {
		return "", "", "", 0, 0, 0, err
	}

	thumbnailDir := filepath.Join(s.baseDir, "thumbnails", userID)
	if err := os.MkdirAll(thumbnailDir, 0o750); err != nil {
		if removeErr := os.Remove(originalFullPath); removeErr != nil {
			log.Printf("failed to remove original file on cleanup: %v", removeErr)
		}
		return "", "", "", 0, 0, 0, err
	}

	thumbnail := resize.Thumbnail(s.thumbnailSize, s.thumbnailSize, img, resize.Lanczos3)
	thumbnailFullPath := filepath.Join(thumbnailDir, imageID+"_thumb.webp")
	thumbFile, err := os.Create(thumbnailFullPath)
	if err != nil {
		if removeErr := os.Remove(originalFullPath); removeErr != nil {
			log.Printf("failed to remove original file on cleanup: %v", removeErr)
		}
		return "", "", "", 0, 0, 0, err
	}

	if err := webp.Encode(thumbFile, thumbnail, &webp.Options{Lossless: false, Quality: 80}); err != nil {
		if closeErr := thumbFile.Close(); closeErr != nil {
			log.Printf("failed to close thumbnail file: %v", closeErr)
		}
		if removeErr := os.Remove(originalFullPath); removeErr != nil {
			log.Printf("failed to remove original file: %v", removeErr)
		}
		if removeErr := os.Remove(thumbnailFullPath); removeErr != nil {
			log.Printf("failed to remove thumbnail file: %v", removeErr)
		}
		return "", "", "", 0, 0, 0, err
	}

	if err := thumbFile.Close(); err != nil {
		log.Printf("failed to close thumbnail file: %v", err)
	}

	originalPath = filepath.Join("originals", userID, imageID+ext)
	thumbnailPath = filepath.Join("thumbnails", userID, imageID+"_thumb.webp")
	return imageID, originalPath, thumbnailPath, width, height, size, nil
}

func (s *Storage) DeleteImage(originalPath, thumbnailPath string) error {
	if err := s.removeFile(filepath.Join(s.baseDir, originalPath)); err != nil {
		return fmt.Errorf("failed to delete original: %w", err)
	}
	if err := s.removeFile(filepath.Join(s.baseDir, thumbnailPath)); err != nil {
		return fmt.Errorf("failed to delete thumbnail: %w", err)
	}
	return nil
}

func (s *Storage) removeFile(filePath string) error {
	if err := s.validatePath(filePath); err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *Storage) validatePath(filePath string) error {
	baseDirClean := filepath.Clean(s.baseDir)
	filePathClean := filepath.Clean(filePath)
	if !strings.HasPrefix(filePathClean, baseDirClean) {
		return fmt.Errorf("path traversal detected: %s", filePath)
	}
	return nil
}

func (s *Storage) GetImagePath(imagePath string) string {
	return filepath.Join(s.baseDir, imagePath)
}
