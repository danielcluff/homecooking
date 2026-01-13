package models

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                uuid.UUID  `json:"id"`
	Title             string     `json:"title"`
	Slug              string     `json:"slug"`
	MarkdownContent   string     `json:"markdown_content"`
	AuthorID          *uuid.UUID `json:"author_id"`
	CategoryID        *uuid.UUID `json:"category_id"`
	Description       *string    `json:"description"`
	PrepTimeMinutes   *int32     `json:"prep_time_minutes"`
	CookTimeMinutes   *int32     `json:"cook_time_minutes"`
	Servings          *int32     `json:"servings"`
	Difficulty        *string    `json:"difficulty"`
	FeaturedImagePath *string    `json:"featured_image_path"`
	IsPublished       bool       `json:"is_published"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	PublishedAt       *time.Time `json:"published_at"`
}

type CreateRecipeRequest struct {
	Title             string  `json:"title"`
	MarkdownContent   string  `json:"markdown_content"`
	CategoryID        *string `json:"category_id"`
	Description       *string `json:"description"`
	PrepTimeMinutes   *int32  `json:"prep_time_minutes"`
	CookTimeMinutes   *int32  `json:"cook_time_minutes"`
	Servings          *int32  `json:"servings"`
	Difficulty        *string `json:"difficulty"`
	FeaturedImagePath *string `json:"featured_image_path"`
	IsPublished       bool    `json:"is_published"`
}

type UpdateRecipeRequest struct {
	Title             *string `json:"title"`
	MarkdownContent   *string `json:"markdown_content"`
	CategoryID        *string `json:"category_id"`
	Description       *string `json:"description"`
	PrepTimeMinutes   *int32  `json:"prep_time_minutes"`
	CookTimeMinutes   *int32  `json:"cook_time_minutes"`
	Servings          *int32  `json:"servings"`
	Difficulty        *string `json:"difficulty"`
	FeaturedImagePath *string `json:"featured_image_path"`
	IsPublished       *bool   `json:"is_published"`
}

type RecipeWithImages struct {
	Recipe     Recipe        `json:"recipe"`
	Category   *Category     `json:"category"`
	Author     *User         `json:"author"`
	Tags       []Tag         `json:"tags"`
	BodyImages []RecipeImage `json:"body_images"`
	Groups     []RecipeGroup `json:"groups"`
}

type RecipeImage struct {
	ID            uuid.UUID `json:"id"`
	RecipeID      uuid.UUID `json:"recipe_id"`
	FilePath      string    `json:"file_path"`
	WebPPath      *string   `json:"webp_path"`
	ThumbnailPath *string   `json:"thumbnail_path"`
	Caption       *string   `json:"caption"`
	OrderIndex    int       `json:"order_index"`
	UploadedAt    time.Time `json:"uploaded_at"`
}
