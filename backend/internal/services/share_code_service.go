package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/models"
	"github.com/homecooking/backend/internal/repository"
)

type ShareCodeService struct {
	shareCodeRepo *repository.ShareCodeRepository
	recipeRepo    *repository.RecipeRepository
}

func NewShareCodeService(shareCodeRepo *repository.ShareCodeRepository, recipeRepo *repository.RecipeRepository) *ShareCodeService {
	return &ShareCodeService{
		shareCodeRepo: shareCodeRepo,
		recipeRepo:    recipeRepo,
	}
}

func (s *ShareCodeService) CreateShareCode(recipeID string, expiresAt *time.Time, maxUses *int) (*models.ShareCode, error) {
	if recipeID == "" {
		return nil, errors.New("recipe_id is required")
	}

	recipe, err := s.recipeRepo.GetByID(recipeID)
	if err != nil {
		return nil, errors.New("recipe not found")
	}

	if !recipe.IsPublished {
		return nil, errors.New("can only share published recipes")
	}

	code, err := generateShareCode()
	if err != nil {
		return nil, err
	}

	shareCode := &models.ShareCode{
		RecipeID:  uuid.MustParse(recipeID),
		Code:      code,
		ExpiresAt: expiresAt,
		MaxUses:   maxUses,
		UseCount:  0,
	}

	return s.shareCodeRepo.Create(shareCode)
}

func (s *ShareCodeService) GetShareCode(code string) (*models.ShareCodeWithRecipe, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	shareCode, err := s.shareCodeRepo.GetByCode(code)
	if err != nil {
		return nil, errors.New("invalid or expired share code")
	}

	if shareCode.MaxUses != nil && shareCode.UseCount >= *shareCode.MaxUses {
		return nil, errors.New("share code has reached maximum uses")
	}

	if shareCode.ExpiresAt != nil && time.Now().After(*shareCode.ExpiresAt) {
		return nil, errors.New("share code has expired")
	}

	return shareCode, nil
}

func (s *ShareCodeService) UseShareCode(code string) error {
	if code == "" {
		return errors.New("code is required")
	}

	shareCode, err := s.shareCodeRepo.GetByCode(code)
	if err != nil {
		return errors.New("invalid or expired share code")
	}

	if shareCode.MaxUses != nil && shareCode.UseCount >= *shareCode.MaxUses {
		return errors.New("share code has reached maximum uses")
	}

	if shareCode.ExpiresAt != nil && time.Now().After(*shareCode.ExpiresAt) {
		return errors.New("share code has expired")
	}

	return s.shareCodeRepo.IncrementUse(shareCode.ID.String())
}

func (s *ShareCodeService) GetShareCodesForRecipe(recipeID string) ([]*models.ShareCode, error) {
	if recipeID == "" {
		return nil, errors.New("recipe_id is required")
	}

	return s.shareCodeRepo.GetForRecipe(recipeID)
}

func (s *ShareCodeService) DeleteShareCode(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return s.shareCodeRepo.Delete(id)
}

func generateShareCode() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
