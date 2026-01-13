package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type TagRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewTagRepository(db *sql.DB, q *sqlc.Queries) *TagRepository {
	return &TagRepository{
		db: db,
		q:  q,
	}
}

func (r *TagRepository) Create(tag *models.Tag) (*models.Tag, error) {
	ctx := context.Background()

	id := tag.ID
	if (uuid.UUID{}) == id {
		id = uuid.New()
	}

	result, err := r.q.CreateTag(ctx, sqlc.CreateTagParams{
		ID:    id,
		Name:  tag.Name,
		Slug:  tag.Slug,
		Color: sql.NullString{String: tag.Color, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *TagRepository) GetByID(id string) (*models.Tag, error) {
	ctx := context.Background()
	result, err := r.q.GetTagByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *TagRepository) GetBySlug(slug string) (*models.Tag, error) {
	ctx := context.Background()
	result, err := r.q.GetTagBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *TagRepository) List() ([]*models.Tag, error) {
	ctx := context.Background()
	results, err := r.q.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]*models.Tag, len(results))
	for i, result := range results {
		tags[i] = r.sqlcToModel(result)
	}
	return tags, nil
}

func (r *TagRepository) Update(id string, tag *models.Tag) (*models.Tag, error) {
	ctx := context.Background()
	result, err := r.q.UpdateTag(ctx, sqlc.UpdateTagParams{
		ID:    uuid.MustParse(id),
		Name:  sql.NullString{String: tag.Name, Valid: true},
		Slug:  sql.NullString{String: tag.Slug, Valid: true},
		Color: sql.NullString{String: tag.Color, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *TagRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteTag(ctx, uuid.MustParse(id))
}

func (r *TagRepository) AddToRecipe(recipeID string, tagID string) error {
	ctx := context.Background()
	return r.q.AddTagToRecipe(ctx, sqlc.AddTagToRecipeParams{
		RecipeID: uuid.MustParse(recipeID),
		TagID:    uuid.MustParse(tagID),
	})
}

func (r *TagRepository) RemoveFromRecipe(recipeID string, tagID string) error {
	ctx := context.Background()
	return r.q.RemoveTagFromRecipe(ctx, sqlc.RemoveTagFromRecipeParams{
		RecipeID: uuid.MustParse(recipeID),
		TagID:    uuid.MustParse(tagID),
	})
}

func (r *TagRepository) GetRecipeTags(recipeID string) ([]*models.Tag, error) {
	ctx := context.Background()
	results, err := r.q.GetRecipeTags(ctx, uuid.MustParse(recipeID))
	if err != nil {
		return nil, err
	}

	tags := make([]*models.Tag, len(results))
	for i, result := range results {
		tags[i] = r.sqlcToModel(result)
	}
	return tags, nil
}

func (r *TagRepository) sqlcToModel(dbTag sqlc.Tag) *models.Tag {
	color := dbTag.Color.String
	if !dbTag.Color.Valid {
		color = "#6366f1"
	}
	return &models.Tag{
		ID:        dbTag.ID,
		Name:      dbTag.Name,
		Slug:      dbTag.Slug,
		Color:     color,
		CreatedAt: dbTag.CreatedAt.Time,
	}
}
