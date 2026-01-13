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

type UserInviteService struct {
	inviteRepo *repository.UserInviteRepository
	userRepo   *repository.UserRepository
}

func NewUserInviteService(inviteRepo *repository.UserInviteRepository, userRepo *repository.UserRepository) *UserInviteService {
	return &UserInviteService{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
	}
}

func (s *UserInviteService) CreateInvite(email string, role string, createdBy string, expiresAt *time.Time) (*models.UserInvite, error) {
	if createdBy == "" {
		return nil, errors.New("created_by is required")
	}

	if role == "" {
		role = "user"
	}

	if role != "user" && role != "admin" {
		return nil, errors.New("invalid role")
	}

	code, err := generateInviteCode()
	if err != nil {
		return nil, err
	}

	createdByUUID := uuid.MustParse(createdBy)

	invite := &models.UserInvite{
		Code:      code,
		Email:     &email,
		Role:      role,
		CreatedBy: &createdByUUID,
		ExpiresAt: expiresAt,
	}

	return s.inviteRepo.Create(invite)
}

func (s *UserInviteService) GetInvite(code string) (*models.UserInvite, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	invite, err := s.inviteRepo.GetByCode(code)
	if err != nil {
		return nil, errors.New("invalid or expired invite code")
	}

	if invite.UsedAt != nil {
		return nil, errors.New("invite has already been used")
	}

	if invite.ExpiresAt != nil && time.Now().After(*invite.ExpiresAt) {
		return nil, errors.New("invite has expired")
	}

	return invite, nil
}

func (s *UserInviteService) UseInvite(code string, usedBy string) (*models.UserInvite, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}

	if usedBy == "" {
		return nil, errors.New("used_by is required")
	}

	invite, err := s.GetInvite(code)
	if err != nil {
		return nil, err
	}

	if invite.Email != nil {
		user, err := s.userRepo.GetByEmail(*invite.Email)
		if err == nil && user != nil {
			return nil, errors.New("user with this email already exists")
		}
	}

	return s.inviteRepo.Use(invite.ID.String(), usedBy)
}

func (s *UserInviteService) ListInvites() ([]*models.UserInvite, error) {
	return s.inviteRepo.List()
}

func (s *UserInviteService) DeleteInvite(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return s.inviteRepo.Delete(id)
}

func generateInviteCode() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
