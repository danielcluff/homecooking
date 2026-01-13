package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/homecooking/backend/internal/services"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

type ExtractFromImageRequest struct {
	ImageType string `json:"image_type"`
}

func (h *AIHandler) ExtractFromImage(w http.ResponseWriter, r *http.Request) {
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusBadRequest)
		return
	}

	imageType := r.URL.Query().Get("type")
	if imageType == "" {
		imageType = "image/jpeg"
	}

	req := &services.ExtractRecipeRequest{
		ImageData: imageData,
		ImageType: imageType,
	}

	result, err := h.aiService.ExtractRecipe(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type EnhanceRecipeRequest struct {
	CurrentTitle       string `json:"current_title"`
	CurrentContent     string `json:"current_content"`
	EnhanceDescription bool   `json:"enhance_description"`
	ImproveStructure   bool   `json:"improve_structure"`
	AddTips            bool   `json:"add_tips"`
}

func (h *AIHandler) EnhanceRecipe(w http.ResponseWriter, r *http.Request) {
	var req EnhanceRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	serviceReq := &services.EnhanceRecipeRequest{
		CurrentTitle:       req.CurrentTitle,
		CurrentContent:     req.CurrentContent,
		EnhanceDescription: req.EnhanceDescription || true,
		ImproveStructure:   req.ImproveStructure || true,
		AddTips:            req.AddTips || true,
	}

	result, err := h.aiService.EnhanceRecipe(serviceReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AIHandler) CheckEnabled(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"enabled": h.aiService.IsEnabled(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AIHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.aiService.GetConfig()

	response := map[string]interface{}{
		"enabled":  cfg.Enabled,
		"provider": cfg.Provider,
		"model":    cfg.Model,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
