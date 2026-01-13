package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/services"
)

type RecipeGroupHandler struct {
	service *services.RecipeGroupService
}

func NewRecipeGroupHandler(service *services.RecipeGroupService) *RecipeGroupHandler {
	return &RecipeGroupHandler{
		service: service,
	}
}

func (h *RecipeGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"groups": groups,
	})
}

func (h *RecipeGroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

func (h *RecipeGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	group := &models.RecipeGroup{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	createdGroup, err := h.service.Create(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

func (h *RecipeGroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var group models.RecipeGroup

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedGroup, err := h.service.Update(id, &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGroup)
}

func (h *RecipeGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (h *RecipeGroupHandler) GetGroupRecipes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	recipes, err := h.service.GetRecipesInGroup(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recipes": recipes,
	})
}

func (h *RecipeGroupHandler) AddRecipeToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	var req struct {
		RecipeID string `json:"recipe_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.RecipeID == "" {
		http.Error(w, "recipe_id is required", http.StatusBadRequest)
		return
	}

	if err := h.service.AddRecipeToGroup(groupID, req.RecipeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (h *RecipeGroupHandler) RemoveRecipeFromGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("id")
	recipeID := r.PathValue("recipeId")

	if err := h.service.RemoveRecipeFromGroup(groupID, recipeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
