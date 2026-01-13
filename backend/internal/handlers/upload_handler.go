package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/homecooking/backend/internal/services"
)

type UploadHandler struct {
	storage *services.StorageService
}

func NewUploadHandler(storage *services.StorageService) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

type UploadResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename,omitempty"`
	URL      string `json:"url,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("Failed to parse form: %v", err)
		sendError(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		sendError(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size == 0 {
		sendError(w, "Empty file", http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		buffer := make([]byte, 512)
		file.Read(buffer)
		file.Seek(0, 0)
		contentType = http.DetectContentType(buffer)
	}

	if !isValidImageType(contentType) {
		sendError(w, "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed", http.StatusBadRequest)
		return
	}

	filename, err := h.storage.SaveImage(header, "recipe")
	if err != nil {
		log.Printf("Failed to save image: %v", err)
		sendError(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	response := UploadResponse{
		Success:  true,
		Filename: filename,
		URL:      fmt.Sprintf("/uploads/%s", filename),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(UploadResponse{
		Success: false,
		Error:   message,
	})
}
