package services

import (
	"errors"
	"strings"

	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	if category.Name == "" {
		return nil, errors.New("name is required")
	}
	if category.Slug == "" {
		category.Slug = generateCategorySlug(category.Name)
	}

	return s.categoryRepo.Create(category)
}

func (s *CategoryService) GetCategory(id string) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) GetBySlug(slug string) (*models.Category, error) {
	return s.categoryRepo.GetBySlug(slug)
}

func (s *CategoryService) ListCategories() ([]*models.Category, error) {
	return s.categoryRepo.List()
}

func (s *CategoryService) UpdateCategory(id string, category *models.Category) (*models.Category, error) {
	if category.Name == "" {
		return nil, errors.New("name is required")
	}
	if category.Slug == "" && category.Name != "" {
		category.Slug = generateCategorySlug(category.Name)
	}

	return s.categoryRepo.Update(id, category)
}

func (s *CategoryService) DeleteCategory(id string) error {
	return s.categoryRepo.Delete(id)
}

func generateCategorySlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")
	return slug
}
