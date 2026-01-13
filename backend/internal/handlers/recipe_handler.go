package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/services"
)

type RecipeHandler struct {
	recipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

func (h *RecipeHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
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

	recipes, err := h.recipeService.ListRecipes(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	recipe, err := h.recipeService.GetRecipe(id)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) GetRecipeBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "Slug required", http.StatusBadRequest)
		return
	}

	recipe, err := h.recipeService.GetRecipeBySlug(slug)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	recipe, err := h.recipeService.CreateRecipe(&req, user.ID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	var req models.UpdateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	recipe, err := h.recipeService.UpdateRecipe(id, &req, user.ID.String())
	if err != nil {
		if err.Error() == "unauthorized: you can only edit your own recipes" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	err := h.recipeService.DeleteRecipe(id, user.ID.String())
	if err != nil {
		if err.Error() == "unauthorized: you can only delete your own recipes" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) PublishRecipe(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	var req struct {
		Published bool `json:"published"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(*models.User)

	recipe, err := h.recipeService.PublishRecipe(id, user.ID.String(), req.Published)
	if err != nil {
		if err.Error() == "unauthorized: you can only publish your own recipes" {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Recipe not found", http.StatusNotFound)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query required", http.StatusBadRequest)
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

	recipes, err := h.recipeService.SearchRecipes(query, limit, offset)
	if err != nil {
		http.Error(w, "Failed to search recipes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
