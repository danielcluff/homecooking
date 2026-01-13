package services

import (
	"errors"
	"strings"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

func (s *TagService) CreateTag(tag *models.Tag) (*models.Tag, error) {
	if tag.Name == "" {
		return nil, errors.New("name is required")
	}
	if tag.Slug == "" {
		tag.Slug = generateTagSlug(tag.Name)
	}
	if tag.Color == "" {
		tag.Color = "#6366f1"
	}

	return s.tagRepo.Create(tag)
}

func (s *TagService) GetTag(id string) (*models.Tag, error) {
	return s.tagRepo.GetByID(id)
}

func (s *TagService) GetBySlug(slug string) (*models.Tag, error) {
	return s.tagRepo.GetBySlug(slug)
}

func (s *TagService) ListTags() ([]*models.Tag, error) {
	return s.tagRepo.List()
}

func (s *TagService) UpdateTag(id string, tag *models.Tag) (*models.Tag, error) {
	if tag.Name == "" {
		return nil, errors.New("name is required")
	}
	if tag.Slug == "" && tag.Name != "" {
		tag.Slug = generateTagSlug(tag.Name)
	}
	if tag.Color == "" {
		tag.Color = "#6366f1"
	}

	return s.tagRepo.Update(id, tag)
}

func (s *TagService) DeleteTag(id string) error {
	return s.tagRepo.Delete(id)
}

func (s *TagService) GetRecipeTags(recipeID string) ([]*models.Tag, error) {
	return s.tagRepo.GetRecipeTags(recipeID)
}

func (s *TagService) AddTagToRecipe(recipeID string, tagID string) error {
	return s.tagRepo.AddToRecipe(recipeID, tagID)
}

func (s *TagService) RemoveTagFromRecipe(recipeID string, tagID string) error {
	return s.tagRepo.RemoveFromRecipe(recipeID, tagID)
}

func generateTagSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}
