package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type RecipeRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewRecipeRepository(db *sql.DB, q *sqlc.Queries) *RecipeRepository {
	return &RecipeRepository{
		db: db,
		q:  q,
	}
}

func (r *RecipeRepository) Create(recipe *models.Recipe) (*models.Recipe, error) {
	ctx := context.Background()

	id := recipe.ID
	if (uuid.UUID{}) == id {
		id = uuid.New()
	}

	result, err := r.q.CreateRecipe(ctx, sqlc.CreateRecipeParams{
		ID:                id,
		Title:             recipe.Title,
		Slug:              recipe.Slug,
		MarkdownContent:   recipe.MarkdownContent,
		AuthorID:          sqlNullUUID(recipe.AuthorID),
		CategoryID:        sqlNullUUID(recipe.CategoryID),
		Description:       sqlNullString(recipe.Description),
		PrepTimeMinutes:   sqlNullInt32(recipe.PrepTimeMinutes),
		CookTimeMinutes:   sqlNullInt32(recipe.CookTimeMinutes),
		Servings:          sqlNullInt32(recipe.Servings),
		Difficulty:        sqlNullString(recipe.Difficulty),
		FeaturedImagePath: sqlNullString(recipe.FeaturedImagePath),
		IsPublished:       sql.NullBool{Bool: recipe.IsPublished, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *RecipeRepository) GetByID(id string) (*models.Recipe, error) {
	ctx := context.Background()
	result, err := r.q.GetRecipeByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *RecipeRepository) GetBySlug(slug string) (*models.Recipe, error) {
	ctx := context.Background()
	result, err := r.q.GetRecipeBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *RecipeRepository) List(limit, offset int) ([]*models.Recipe, error) {
	ctx := context.Background()
	results, err := r.q.ListRecipes(ctx, sqlc.ListRecipesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	recipes := make([]*models.Recipe, len(results))
	for i, result := range results {
		recipes[i] = r.sqlcToModel(result)
	}
	return recipes, nil
}

func (r *RecipeRepository) Search(query string, limit, offset int) ([]*models.Recipe, error) {
	ctx := context.Background()
	results, err := r.q.SearchRecipes(ctx, sqlc.SearchRecipesParams{
		Column1: sql.NullString{String: query, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	recipes := make([]*models.Recipe, len(results))
	for i, result := range results {
		recipes[i] = r.sqlcToModel(result)
	}
	return recipes, nil
}

func (r *RecipeRepository) Update(id string, recipe *models.Recipe) (*models.Recipe, error) {
	ctx := context.Background()
	result, err := r.q.UpdateRecipe(ctx, sqlc.UpdateRecipeParams{
		ID:                uuid.MustParse(id),
		Title:             sqlNullStringPtr(recipe.Title),
		MarkdownContent:   sqlNullStringPtr(recipe.MarkdownContent),
		CategoryID:        sqlNullUUID(recipe.CategoryID),
		Description:       sqlNullString(recipe.Description),
		PrepTimeMinutes:   sqlNullInt32(recipe.PrepTimeMinutes),
		CookTimeMinutes:   sqlNullInt32(recipe.CookTimeMinutes),
		Servings:          sqlNullInt32(recipe.Servings),
		Difficulty:        sqlNullString(recipe.Difficulty),
		FeaturedImagePath: sqlNullString(recipe.FeaturedImagePath),
		IsPublished:       sqlNullBoolPtr(recipe.IsPublished),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *RecipeRepository) UpdatePublishedStatus(id string, isPublished bool) (*models.Recipe, error) {
	ctx := context.Background()
	result, err := r.q.UpdateRecipePublishedStatus(ctx, sqlc.UpdateRecipePublishedStatusParams{
		ID:          uuid.MustParse(id),
		IsPublished: sql.NullBool{Bool: isPublished, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *RecipeRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteRecipe(ctx, uuid.MustParse(id))
}

func (r *RecipeRepository) UpdateFeaturedImage(id string, imagePath *string) error {
	ctx := context.Background()
	_, err := r.q.UpdateRecipeFeaturedImage(ctx, sqlc.UpdateRecipeFeaturedImageParams{
		ID:                uuid.MustParse(id),
		FeaturedImagePath: sqlNullString(imagePath),
	})
	return err
}

func (r *RecipeRepository) sqlcToModel(dbRecipe sqlc.Recipe) *models.Recipe {
	return &models.Recipe{
		ID:                dbRecipe.ID,
		Title:             dbRecipe.Title,
		Slug:              dbRecipe.Slug,
		MarkdownContent:   dbRecipe.MarkdownContent,
		AuthorID:          nullUUIDToPtr(dbRecipe.AuthorID),
		CategoryID:        nullUUIDToPtr(dbRecipe.CategoryID),
		Description:       nullStringToPtr(dbRecipe.Description),
		PrepTimeMinutes:   nullInt32ToPtr(dbRecipe.PrepTimeMinutes),
		CookTimeMinutes:   nullInt32ToPtr(dbRecipe.CookTimeMinutes),
		Servings:          nullInt32ToPtr(dbRecipe.Servings),
		Difficulty:        nullStringToPtr(dbRecipe.Difficulty),
		FeaturedImagePath: nullStringToPtr(dbRecipe.FeaturedImagePath),
		IsPublished:       dbRecipe.IsPublished.Bool,
		CreatedAt:         dbRecipe.CreatedAt.Time,
		UpdatedAt:         dbRecipe.UpdatedAt.Time,
		PublishedAt:       nullTimeToTimePtr(dbRecipe.PublishedAt),
	}
}
