package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/services"
)

type TagHandler struct {
	tagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagService.ListTags()
	if err != nil {
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Tag ID required", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.GetTag(id)
	if err != nil {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req models.Tag
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.CreateTag(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Tag ID required", http.StatusBadRequest)
		return
	}

	var req models.Tag
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := h.tagService.UpdateTag(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Tag ID required", http.StatusBadRequest)
		return
	}

	err := h.tagService.DeleteTag(id)
	if err != nil {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TagHandler) GetRecipeTags(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("recipeId")
	if recipeID == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}

	tags, err := h.tagService.GetRecipeTags(recipeID)
	if err != nil {
		http.Error(w, "Failed to fetch recipe tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) AddTagToRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("recipeId")
	tagID := r.PathValue("tagId")

	if recipeID == "" || tagID == "" {
		http.Error(w, "Recipe ID and Tag ID required", http.StatusBadRequest)
		return
	}

	err := h.tagService.AddTagToRecipe(recipeID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TagHandler) RemoveTagFromRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := r.PathValue("recipeId")
	tagID := r.PathValue("tagId")

	if recipeID == "" || tagID == "" {
		http.Error(w, "Recipe ID and Tag ID required", http.StatusBadRequest)
		return
	}

	err := h.tagService.RemoveTagFromRecipe(recipeID, tagID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
