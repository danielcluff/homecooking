package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/services"
)

type VariationHandler struct {
	variationService *services.VariationService
}

func NewVariationHandler(variationService *services.VariationService) *VariationHandler {
	return &VariationHandler{
		variationService: variationService,
	}
}

func (h *VariationHandler) ListVariations(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("id")
	if recipeID == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	variations, err := h.variationService.GetVariationsByRecipe(recipeID)
	if err != nil {
		http.Error(w, "Failed to fetch variations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variations)
}

func (h *VariationHandler) GetVariation(w http.ResponseWriter, r *http.Request) {
	variationID := r.PathValue("variationId")
	if variationID == "" {
		http.Error(w, "Variation ID required", http.StatusBadRequest)
		return
	}

	variation, err := h.variationService.GetVariation(variationID)
	if err != nil {
		http.Error(w, "Variation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variation)
}

func (h *VariationHandler) CreateVariation(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("id")
	if recipeID == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	var req models.CreateVariationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	variation, err := h.variationService.CreateVariation(&req, recipeID, user.ID.String())
	if err != nil {
		if err.Error() == "markdown content is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to create variation", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(variation)
}

func (h *VariationHandler) UpdateVariation(w http.ResponseWriter, r *http.Request) {
	variationID := r.PathValue("variationId")
	if variationID == "" {
		http.Error(w, "Variation ID required", http.StatusBadRequest)
		return
	}

	var req models.UpdateVariationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	variation, err := h.variationService.UpdateVariation(variationID, &req, user.ID.String())
	if err != nil {
		if err.Error() == "unauthorized: you can only edit your own variations" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Failed to update variation", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variation)
}

func (h *VariationHandler) DeleteVariation(w http.ResponseWriter, r *http.Request) {
	variationID := r.PathValue("variationId")
	if variationID == "" {
		http.Error(w, "Variation ID required", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	err := h.variationService.DeleteVariation(variationID, user.ID.String())
	if err != nil {
		if err.Error() == "unauthorized: you can only delete your own variations" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Failed to delete variation", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VariationHandler) GetVariationsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorID := r.PathValue("authorId")
	if authorID == "" {
		http.Error(w, "Author ID required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	variations, err := h.variationService.ListVariationsByAuthor(authorID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch variations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variations)
}
