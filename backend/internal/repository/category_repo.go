package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewCategoryRepository(db *sql.DB, q *sqlc.Queries) *CategoryRepository {
	return &CategoryRepository{
		db: db,
		q:  q,
	}
}

func (r *CategoryRepository) Create(category *models.Category) (*models.Category, error) {
	ctx := context.Background()

	id := category.ID
	if (uuid.UUID{}) == id {
		id = uuid.New()
	}

	result, err := r.q.CreateCategory(ctx, sqlc.CreateCategoryParams{
		ID:          id,
		Name:        category.Name,
		Slug:        category.Slug,
		Icon:        sqlNullString(category.Icon),
		Description: sqlNullString(category.Description),
		OrderIndex:  sql.NullInt32{Int32: int32(category.OrderIndex), Valid: category.OrderIndex != 0},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *CategoryRepository) GetByID(id string) (*models.Category, error) {
	ctx := context.Background()
	result, err := r.q.GetCategoryByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	ctx := context.Background()
	result, err := r.q.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *CategoryRepository) List() ([]*models.Category, error) {
	ctx := context.Background()
	results, err := r.q.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]*models.Category, len(results))
	for i, result := range results {
		categories[i] = r.sqlcToModel(result)
	}
	return categories, nil
}

func (r *CategoryRepository) Update(id string, category *models.Category) (*models.Category, error) {
	ctx := context.Background()
	result, err := r.q.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:          uuid.MustParse(id),
		Name:        sql.NullString{String: category.Name, Valid: true},
		Slug:        sql.NullString{String: category.Slug, Valid: true},
		Icon:        sqlNullString(category.Icon),
		Description: sqlNullString(category.Description),
		OrderIndex:  sql.NullInt32{Int32: int32(category.OrderIndex), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *CategoryRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteCategory(ctx, uuid.MustParse(id))
}

func (r *CategoryRepository) sqlcToModel(dbCategory sqlc.Category) *models.Category {
	orderIndex := int(dbCategory.OrderIndex.Int32)
	return &models.Category{
		ID:          dbCategory.ID,
		Name:        dbCategory.Name,
		Slug:        dbCategory.Slug,
		Icon:        nullStringToPtr(dbCategory.Icon),
		Description: nullStringToPtr(dbCategory.Description),
		OrderIndex:  orderIndex,
		CreatedAt:   dbCategory.CreatedAt.Time,
	}
}
