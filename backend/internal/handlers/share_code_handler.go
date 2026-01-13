package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/services"
)

type ShareCodeHandler struct {
	shareCodeService *services.ShareCodeService
}

func NewShareCodeHandler(shareCodeService *services.ShareCodeService) *ShareCodeHandler {
	return &ShareCodeHandler{
		shareCodeService: shareCodeService,
	}
}

type CreateShareCodeRequest struct {
	RecipeID  string     `json:"recipe_id"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   *int       `json:"max_uses"`
}

func (h *ShareCodeHandler) CreateShareCode(w http.ResponseWriter, r *http.Request) {
	var req CreateShareCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shareCode, err := h.shareCodeService.CreateShareCode(req.RecipeID, req.ExpiresAt, req.MaxUses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shareCode)
}

func (h *ShareCodeHandler) GetShareCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Code required", http.StatusBadRequest)
		return
	}

	shareCode, err := h.shareCodeService.GetShareCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shareCode)
}

func (h *ShareCodeHandler) GetShareCodesForRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("recipeId")
	if recipeID == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	shareCodes, err := h.shareCodeService.GetShareCodesForRecipe(recipeID)
	if err != nil {
		http.Error(w, "Failed to fetch share codes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shareCodes)
}

func (h *ShareCodeHandler) DeleteShareCode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	if err := h.shareCodeService.DeleteShareCode(id); err != nil {
		http.Error(w, "Failed to delete share code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AccessRecipeByShareCodeResponse struct {
	ID                uuid.UUID  `json:"id"`
	Title             string     `json:"title"`
	Slug              string     `json:"slug"`
	MarkdownContent   string     `json:"markdown_content"`
	Description       *string    `json:"description"`
	PrepTimeMinutes   *int       `json:"prep_time_minutes"`
	CookTimeMinutes   *int       `json:"cook_time_minutes"`
	Servings          *int       `json:"servings"`
	Difficulty        *string    `json:"difficulty"`
	FeaturedImagePath *string    `json:"featured_image_path"`
	PublishedAt       *time.Time `json:"published_at"`
	CategoryID        *uuid.UUID `json:"category_id"`
}

func (h *ShareCodeHandler) AccessRecipeByShareCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Code required", http.StatusBadRequest)
		return
	}

	shareCode, err := h.shareCodeService.GetShareCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := h.shareCodeService.UseShareCode(code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shareCode)
}
