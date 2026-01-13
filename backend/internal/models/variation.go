package models

import (
	"time"

	"github.com/google/uuid"
)

type RecipeVariation struct {
	ID              uuid.UUID `json:"id"`
	RecipeID        uuid.UUID `json:"recipe_id"`
	AuthorID        uuid.UUID `json:"author_id"`
	MarkdownContent string    `json:"markdown_content"`
	PrepTimeMinutes *int32    `json:"prep_time_minutes"`
	CookTimeMinutes *int32    `json:"cook_time_minutes"`
	Servings        *int32    `json:"servings"`
	Difficulty      *string   `json:"difficulty"`
	Notes           *string   `json:"notes"`
	IsPublished     bool      `json:"is_published"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateVariationRequest struct {
	MarkdownContent string  `json:"markdown_content"`
	PrepTimeMinutes *int32  `json:"prep_time_minutes"`
	CookTimeMinutes *int32  `json:"cook_time_minutes"`
	Servings        *int32  `json:"servings"`
	Difficulty      *string `json:"difficulty"`
	Notes           *string `json:"notes"`
	IsPublished     bool    `json:"is_published"`
}

type UpdateVariationRequest struct {
	MarkdownContent *string `json:"markdown_content"`
	PrepTimeMinutes *int32  `json:"prep_time_minutes"`
	CookTimeMinutes *int32  `json:"cook_time_minutes"`
	Servings        *int32  `json:"servings"`
	Difficulty      *string `json:"difficulty"`
	Notes           *string `json:"notes"`
	IsPublished     *bool   `json:"is_published"`
}

type VariationWithAuthor struct {
	RecipeVariation
	Author *User `json:"author"`
}

type RecipeWithVariations struct {
	Recipe     Recipe                `json:"recipe"`
	Category   *Category             `json:"category"`
	Author     *User                 `json:"author"`
	Tags       []Tag                 `json:"tags"`
	BodyImages []RecipeImage         `json:"body_images"`
	Groups     []RecipeGroup         `json:"groups"`
	Variations []VariationWithAuthor `json:"variations"`
}
