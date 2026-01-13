package models

import (
	"time"

	"github.com/google/uuid"
)

type ShareCode struct {
	ID        uuid.UUID  `json:"id"`
	RecipeID  uuid.UUID  `json:"recipe_id"`
	Code      string     `json:"code"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   *int       `json:"max_uses"`
	UseCount  int        `json:"use_count"`
	CreatedAt time.Time  `json:"created_at"`
}

type UserInvite struct {
	ID        uuid.UUID  `json:"id"`
	Code      string     `json:"code"`
	Email     *string    `json:"email"`
	Role      string     `json:"role"`
	CreatedBy *uuid.UUID `json:"created_by"`
	ExpiresAt *time.Time `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	UsedBy    *uuid.UUID `json:"used_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type AppSetting struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShareCodeWithRecipe struct {
	ID          uuid.UUID  `json:"id"`
	RecipeID    uuid.UUID  `json:"recipe_id"`
	Code        string     `json:"code"`
	ExpiresAt   *time.Time `json:"expires_at"`
	MaxUses     *int       `json:"max_uses"`
	UseCount    int        `json:"use_count"`
	CreatedAt   time.Time  `json:"created_at"`
	RecipeTitle string     `json:"recipe_title"`
	RecipeSlug  string     `json:"recipe_slug"`
}
