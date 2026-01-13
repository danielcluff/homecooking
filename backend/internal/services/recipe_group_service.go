package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type RecipeGroupService struct {
	repo *repository.RecipeGroupRepository
}

func NewRecipeGroupService(repo *repository.RecipeGroupRepository) *RecipeGroupService {
	return &RecipeGroupService{
		repo: repo,
	}
}

func (s *RecipeGroupService) Create(group *models.RecipeGroup) (*models.RecipeGroup, error) {
	if group.Slug == "" {
		group.Slug = s.GenerateSlug(group.Name)
	}
	return s.repo.Create(group)
}

func (s *RecipeGroupService) GetByID(id string) (*models.RecipeGroup, error) {
	return s.repo.GetByID(id)
}

func (s *RecipeGroupService) GetBySlug(slug string) (*models.RecipeGroup, error) {
	return s.repo.GetBySlug(slug)
}

func (s *RecipeGroupService) List() ([]*models.RecipeGroup, error) {
	return s.repo.List()
}

func (s *RecipeGroupService) Update(id string, group *models.RecipeGroup) (*models.RecipeGroup, error) {
	if id != "" {
		_, err := s.GetByID(id)
		if err != nil {
			return nil, err
		}
	}

	if group.Slug == "" {
		group.Slug = s.GenerateSlug(group.Name)
	}
	return s.repo.Update(id, group)
}

func (s *RecipeGroupService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *RecipeGroupService) AddRecipeToGroup(groupID, recipeID string) error {
	return s.repo.AddRecipeToGroup(groupID, recipeID)
}

func (s *RecipeGroupService) RemoveRecipeFromGroup(groupID, recipeID string) error {
	return s.repo.RemoveRecipeFromGroup(groupID, recipeID)
}

func (s *RecipeGroupService) GetRecipesInGroup(groupID string) ([]*models.Recipe, error) {
	return s.repo.GetRecipesInGroup(groupID)
}

func (s *RecipeGroupService) GenerateSlug(name string) string {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, slug)

	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("group-%s", uuid.New().String()[:8])
	}
	return slug
}
