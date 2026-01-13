package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type StorageService struct {
	localPath string
	maxSize   int64
}

func NewStorageService(localPath string, maxSize int64) *StorageService {
	return &StorageService{
		localPath: localPath,
		maxSize:   maxSize,
	}
}

func (s *StorageService) SaveImage(file *multipart.FileHeader, prefix string) (string, error) {
	if file.Size > s.maxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	filename := s.generateFilename(prefix, ext)
	dstPath := filepath.Join(s.localPath, filename)

	img, _, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	resized, err := s.resizeImage(img, 1200, 800)
	if err != nil {
		return "", fmt.Errorf("failed to resize image: %w", err)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if ext == ".png" {
		err = png.Encode(dst, resized)
	} else {
		opts := &jpeg.Options{Quality: 85}
		err = jpeg.Encode(dst, resized, opts)
	}
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}

func (s *StorageService) resizeImage(img image.Image, maxWidth, maxHeight int) (image.Image, error) {
	bounds := img.Bounds()

	if bounds.Dx() <= maxWidth && bounds.Dy() <= maxHeight {
		return img, nil
	}

	ratio := float64(maxWidth) / float64(bounds.Dx())
	if float64(bounds.Dy())*ratio > float64(maxHeight) {
		ratio = float64(maxHeight) / float64(bounds.Dy())
	}

	newWidth := int(float64(bounds.Dx()) * ratio)
	newHeight := int(float64(bounds.Dy()) * ratio)

	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / ratio)
			srcY := int(float64(y) / ratio)
			if srcX < bounds.Dx() && srcY < bounds.Dy() {
				thumbnail.Set(x, y, img.At(srcX, srcY))
			}
		}
	}

	return thumbnail, nil
}

func (s *StorageService) DeleteImage(filename string) error {
	path := filepath.Join(s.localPath, filename)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *StorageService) isValidImageExtension(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExts[ext]
}

func (s *StorageService) generateFilename(prefix string, ext string) string {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randomStr := hex.EncodeToString(randBytes)
	return fmt.Sprintf("%s_%s%s", prefix, randomStr, ext)
}

func (s *StorageService) EnsureDirectory() error {
	if err := os.MkdirAll(s.localPath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}
	return nil
}
