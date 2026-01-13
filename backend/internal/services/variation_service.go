package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type VariationService struct {
	variationRepo *repository.VariationRepository
	recipeRepo    *repository.RecipeRepository
}

func NewVariationService(variationRepo *repository.VariationRepository, recipeRepo *repository.RecipeRepository) *VariationService {
	return &VariationService{
		variationRepo: variationRepo,
		recipeRepo:    recipeRepo,
	}
}

func (s *VariationService) CreateVariation(req *models.CreateVariationRequest, recipeID string, authorID string) (*models.RecipeVariation, error) {
	if req.MarkdownContent == "" {
		return nil, errors.New("markdown content is required")
	}

	recipeUUID := uuid.MustParse(recipeID)
	authorUUID := uuid.MustParse(authorID)

	variationModel := &models.RecipeVariation{
		RecipeID:        recipeUUID,
		AuthorID:        authorUUID,
		MarkdownContent: req.MarkdownContent,
		PrepTimeMinutes: req.PrepTimeMinutes,
		CookTimeMinutes: req.CookTimeMinutes,
		Servings:        req.Servings,
		Difficulty:      req.Difficulty,
		Notes:           req.Notes,
		IsPublished:     req.IsPublished,
	}

	return s.variationRepo.Create(variationModel)
}

func (s *VariationService) GetVariation(id string) (*models.RecipeVariation, error) {
	return s.variationRepo.GetByID(id)
}

func (s *VariationService) GetVariationsByRecipe(recipeID string) ([]*models.RecipeVariation, error) {
	return s.variationRepo.GetByRecipe(recipeID)
}

func (s *VariationService) GetPublishedVariationsByRecipe(recipeID string) ([]*models.RecipeVariation, error) {
	return s.variationRepo.GetPublishedByRecipe(recipeID)
}

func (s *VariationService) GetVariationByRecipeAndAuthor(recipeID string, authorID string) (*models.RecipeVariation, error) {
	return s.variationRepo.GetByRecipeAndAuthor(recipeID, authorID)
}

func (s *VariationService) UpdateVariation(id string, req *models.UpdateVariationRequest, authorID string) (*models.RecipeVariation, error) {
	existing, err := s.variationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID.String() != authorID {
		return nil, errors.New("unauthorized: you can only edit your own variations")
	}

	return s.variationRepo.Update(id, req)
}

func (s *VariationService) DeleteVariation(id string, authorID string) error {
	existing, err := s.variationRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existing.AuthorID.String() != authorID {
		return errors.New("unauthorized: you can only delete your own variations")
	}

	return s.variationRepo.Delete(id)
}

func (s *VariationService) ListVariationsByAuthor(authorID string, limit, offset int) ([]*models.RecipeVariation, error) {
	return s.variationRepo.ListByAuthor(authorID, limit, offset)
}
