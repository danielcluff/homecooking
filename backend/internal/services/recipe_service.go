package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type RecipeService struct {
	recipeRepo *repository.RecipeRepository
}

func NewRecipeService(recipeRepo *repository.RecipeRepository) *RecipeService {
	return &RecipeService{
		recipeRepo: recipeRepo,
	}
}

func (s *RecipeService) CreateRecipe(recipe *models.CreateRecipeRequest, authorID string) (*models.Recipe, error) {
	if recipe.Title == "" {
		return nil, errors.New("title is required")
	}
	if recipe.MarkdownContent == "" {
		return nil, errors.New("markdown content is required")
	}

	slug := generateSlug(recipe.Title)
	authorUUID := uuid.MustParse(authorID)

	recipeModel := &models.Recipe{
		Title:             recipe.Title,
		Slug:              slug,
		MarkdownContent:   recipe.MarkdownContent,
		AuthorID:          &authorUUID,
		CategoryID:        parseUUID(recipe.CategoryID),
		Description:       recipe.Description,
		PrepTimeMinutes:   recipe.PrepTimeMinutes,
		CookTimeMinutes:   recipe.CookTimeMinutes,
		Servings:          recipe.Servings,
		Difficulty:        recipe.Difficulty,
		FeaturedImagePath: recipe.FeaturedImagePath,
		IsPublished:       recipe.IsPublished,
	}

	return s.recipeRepo.Create(recipeModel)
}

func (s *RecipeService) GetRecipe(id string) (*models.Recipe, error) {
	return s.recipeRepo.GetByID(id)
}

func (s *RecipeService) GetRecipeBySlug(slug string) (*models.Recipe, error) {
	return s.recipeRepo.GetBySlug(slug)
}

func (s *RecipeService) ListRecipes(limit, offset int) ([]*models.Recipe, error) {
	return s.recipeRepo.List(limit, offset)
}

func (s *RecipeService) SearchRecipes(query string, limit, offset int) ([]*models.Recipe, error) {
	return s.recipeRepo.Search(query, limit, offset)
}

func (s *RecipeService) UpdateRecipe(id string, req *models.UpdateRecipeRequest, authorID string) (*models.Recipe, error) {
	existing, err := s.recipeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID.String() != authorID {
		return nil, errors.New("unauthorized: you can only edit your own recipes")
	}

	recipeModel := &models.Recipe{
		Title:             toString(req.Title, existing.Title),
		MarkdownContent:   toString(req.MarkdownContent, existing.MarkdownContent),
		CategoryID:        parseUUID(req.CategoryID),
		Description:       req.Description,
		PrepTimeMinutes:   req.PrepTimeMinutes,
		CookTimeMinutes:   req.CookTimeMinutes,
		Servings:          req.Servings,
		Difficulty:        req.Difficulty,
		FeaturedImagePath: req.FeaturedImagePath,
		IsPublished:       toBool(req.IsPublished, existing.IsPublished),
	}

	if req.Title != nil {
		recipeModel.Slug = generateSlug(*req.Title)
	}

	return s.recipeRepo.Update(id, recipeModel)
}

func (s *RecipeService) DeleteRecipe(id string, authorID string) error {
	existing, err := s.recipeRepo.GetByID(id)
	if err != nil {
		return err
	}

	if existing.AuthorID.String() != authorID {
		return errors.New("unauthorized: you can only delete your own recipes")
	}

	return s.recipeRepo.Delete(id)
}

func (s *RecipeService) PublishRecipe(id string, authorID string, published bool) (*models.Recipe, error) {
	existing, err := s.recipeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing.AuthorID.String() != authorID {
		return nil, errors.New("unauthorized: you can only publish your own recipes")
	}

	return s.recipeRepo.UpdatePublishedStatus(id, published)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}

func parseUUID(s *string) *uuid.UUID {
	if s == nil {
		return nil
	}
	parsed, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &parsed
}

func toString(s *string, defaultVal string) string {
	if s == nil {
		return defaultVal
	}
	return *s
}

func toBool(b *bool, defaultVal bool) bool {
	if b == nil {
		return defaultVal
	}
	return *b
}
