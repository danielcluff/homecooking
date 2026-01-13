package services

import (
	"image"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsValidImageExtension(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	tests := []struct {
		name     string
		ext      string
		expected bool
	}{
		{"jpg", ".jpg", true},
		{"jpeg", ".jpeg", true},
		{"png", ".png", true},
		{"gif", ".gif", true},
		{"webp", ".webp", true},
		{"bmp", ".bmp", false},
		{"tiff", ".tiff", false},
		{"pdf", ".pdf", false},
		{"txt", ".txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.isValidImageExtension(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateFilename(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	filename := s.generateFilename("recipe", ".jpg")

	assert.Contains(t, filename, "recipe_")
	assert.True(t, strings.HasSuffix(filename, ".jpg"))
	assert.Equal(t, len(filename), len("recipe_")+16+4) // prefix + 16 chars (hex) + .jpg
}

func TestEnsureDirectory(t *testing.T) {
	dir := "./test-uploads-test"
	defer os.RemoveAll(dir)

	s := NewStorageService(dir, 10*1024*1024)

	err := s.EnsureDirectory()
	require.NoError(t, err)

	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestDeleteImage(t *testing.T) {
	dir := "./test-uploads-test"
	defer os.RemoveAll(dir)

	s := NewStorageService(dir, 10*1024*1024)
	s.EnsureDirectory()

	// Create test file
	filename := s.generateFilename("test", ".jpg")
	path := dir + "/" + filename
	file, err := os.Create(path)
	require.NoError(t, err)
	file.Write([]byte("test"))
	file.Close()

	// Verify file exists
	_, err = os.Stat(path)
	require.NoError(t, err)

	// Delete file
	err = s.DeleteImage(filename)
	require.NoError(t, err)

	// Verify file deleted
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestDeleteImage_NotFound(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	err := s.DeleteImage("nonexistent.jpg")
	assert.NoError(t, err)
}

func TestSaveImage_InvalidType(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	fileHeader := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	_, err := s.SaveImage(fileHeader, "test")

	// Should fail when trying to open the file (no actual file exists)
	assert.Error(t, err)
}

func TestSaveImage_TooLarge(t *testing.T) {
	s := NewStorageService("./test-uploads", 1024)

	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     2048,
	}

	_, err := s.SaveImage(fileHeader, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed size")
}

func TestResizeImage_Small(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	result, err := s.resizeImage(img, 1200, 800)

	require.NoError(t, err)
	assert.Same(t, img, result) // Should return same image if already small enough
}

func TestResizeImage_Large(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	img := image.NewRGBA(image.Rect(0, 0, 2000, 1500))
	result, err := s.resizeImage(img, 1200, 800)

	require.NoError(t, err)
	bounds := result.Bounds()
	assert.LessOrEqual(t, bounds.Dx(), 1200)
	assert.LessOrEqual(t, bounds.Dy(), 800)
}

func TestResizeImage_Ratio(t *testing.T) {
	s := NewStorageService("./test-uploads", 10*1024*1024)

	img := image.NewRGBA(image.Rect(0, 0, 2000, 1000)) // 2:1 ratio
	result, err := s.resizeImage(img, 1200, 800)

	require.NoError(t, err)
	bounds := result.Bounds()

	// Should maintain aspect ratio (roughly)
	ratio := float64(bounds.Dx()) / float64(bounds.Dy())
	assert.InDelta(t, 2.0, ratio, 0.1)
}
