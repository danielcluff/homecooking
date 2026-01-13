package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type VariationRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewVariationRepository(db *sql.DB, q *sqlc.Queries) *VariationRepository {
	return &VariationRepository{
		db: db,
		q:  q,
	}
}

func (r *VariationRepository) Create(variation *models.RecipeVariation) (*models.RecipeVariation, error) {
	ctx := context.Background()

	id := variation.ID
	if (uuid.UUID{}) == id {
		id = uuid.New()
	}

	result, err := r.q.CreateVariation(ctx, sqlc.CreateVariationParams{
		ID:              id,
		RecipeID:        variation.RecipeID,
		AuthorID:        variation.AuthorID,
		MarkdownContent: variation.MarkdownContent,
		PrepTimeMinutes: sqlNullInt32(variation.PrepTimeMinutes),
		CookTimeMinutes: sqlNullInt32(variation.CookTimeMinutes),
		Servings:        sqlNullInt32(variation.Servings),
		Difficulty:      sqlNullString(variation.Difficulty),
		Notes:           sqlNullString(variation.Notes),
		IsPublished:     sqlNullBool(variation.IsPublished),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *VariationRepository) GetByID(id string) (*models.RecipeVariation, error) {
	ctx := context.Background()
	result, err := r.q.GetVariationByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *VariationRepository) GetByRecipe(recipeID string) ([]*models.RecipeVariation, error) {
	ctx := context.Background()
	results, err := r.q.GetVariationsByRecipe(ctx, uuid.MustParse(recipeID))
	if err != nil {
		return nil, err
	}

	variations := make([]*models.RecipeVariation, len(results))
	for i, result := range results {
		variations[i] = r.sqlcToModel(result)
	}
	return variations, nil
}

func (r *VariationRepository) GetPublishedByRecipe(recipeID string) ([]*models.RecipeVariation, error) {
	allVariations, err := r.GetByRecipe(recipeID)
	if err != nil {
		return nil, err
	}

	var publishedVariations []*models.RecipeVariation
	for _, v := range allVariations {
		if v.IsPublished {
			publishedVariations = append(publishedVariations, v)
		}
	}
	return publishedVariations, nil
}

func (r *VariationRepository) GetByRecipeAndAuthor(recipeID string, authorID string) (*models.RecipeVariation, error) {
	ctx := context.Background()
	result, err := r.q.GetVariationByRecipeAndAuthor(ctx, uuid.MustParse(recipeID), uuid.MustParse(authorID))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *VariationRepository) Update(id string, req *models.UpdateVariationRequest) (*models.RecipeVariation, error) {
	ctx := context.Background()
	result, err := r.q.UpdateVariation(ctx, sqlc.UpdateVariationParams{
		ID:              uuid.MustParse(id),
		MarkdownContent: sqlNullString(req.MarkdownContent),
		PrepTimeMinutes: sqlNullInt32(req.PrepTimeMinutes),
		CookTimeMinutes: sqlNullInt32(req.CookTimeMinutes),
		Servings:        sqlNullInt32(req.Servings),
		Difficulty:      sqlNullString(req.Difficulty),
		Notes:           sqlNullString(req.Notes),
		IsPublished:     sqlNullBoolFromPtr(req.IsPublished),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *VariationRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteVariation(ctx, uuid.MustParse(id))
}

func (r *VariationRepository) ListByAuthor(authorID string, limit, offset int) ([]*models.RecipeVariation, error) {
	ctx := context.Background()
	results, err := r.q.ListVariationsByAuthor(ctx, sqlc.ListVariationsByAuthorParams{
		AuthorID: uuid.MustParse(authorID),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	variations := make([]*models.RecipeVariation, len(results))
	for i, result := range results {
		variations[i] = r.sqlcToModel(result)
	}
	return variations, nil
}

func (r *VariationRepository) GetVariationsWithAuthor(recipeID string) ([]*models.VariationWithAuthor, error) {
	ctx := context.Background()
	results, err := r.q.GetVariationsByRecipeWithAuthor(ctx, uuid.MustParse(recipeID))
	if err != nil {
		return nil, err
	}

	variations := make([]*models.VariationWithAuthor, len(results))
	for i, result := range results {
		variations[i] = &models.VariationWithAuthor{
			RecipeVariation: models.RecipeVariation{
				ID:              result.ID,
				RecipeID:        result.RecipeID,
				AuthorID:        result.AuthorID,
				MarkdownContent: result.MarkdownContent,
				PrepTimeMinutes: nullInt32ToPtr(result.PrepTimeMinutes),
				CookTimeMinutes: nullInt32ToPtr(result.CookTimeMinutes),
				Servings:        nullInt32ToPtr(result.Servings),
				Difficulty:      nullStringToPtr(result.Difficulty),
				Notes:           nullStringToPtr(result.Notes),
				IsPublished:     result.IsPublished.Bool,
				CreatedAt:       result.CreatedAt.Time,
				UpdatedAt:       result.UpdatedAt.Time,
			},
			Author: &models.User{
				ID:    result.AuthorID,
				Email: result.AuthorEmail,
				Role:  result.AuthorRole,
			},
		}
	}
	return variations, nil
}

func (r *VariationRepository) sqlcToModel(dbVariation sqlc.RecipeVariation) *models.RecipeVariation {
	return &models.RecipeVariation{
		ID:              dbVariation.ID,
		RecipeID:        dbVariation.RecipeID,
		AuthorID:        dbVariation.AuthorID,
		MarkdownContent: dbVariation.MarkdownContent,
		PrepTimeMinutes: nullInt32ToPtr(dbVariation.PrepTimeMinutes),
		CookTimeMinutes: nullInt32ToPtr(dbVariation.CookTimeMinutes),
		Servings:        nullInt32ToPtr(dbVariation.Servings),
		Difficulty:      nullStringToPtr(dbVariation.Difficulty),
		Notes:           nullStringToPtr(dbVariation.Notes),
		IsPublished:     dbVariation.IsPublished.Bool,
		CreatedAt:       dbVariation.CreatedAt.Time,
		UpdatedAt:       dbVariation.UpdatedAt.Time,
	}
}
