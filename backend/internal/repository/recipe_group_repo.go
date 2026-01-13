package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type RecipeGroupRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewRecipeGroupRepository(db *sql.DB, q *sqlc.Queries) *RecipeGroupRepository {
	return &RecipeGroupRepository{
		db: db,
		q:  q,
	}
}

func (r *RecipeGroupRepository) Create(group *models.RecipeGroup) (*models.RecipeGroup, error) {
	ctx := context.Background()

	id := group.ID
	if (uuid.UUID{}) == id {
		id = uuid.New()
	}

	result, err := r.q.CreateRecipeGroup(ctx, sqlc.CreateRecipeGroupParams{
		ID:          id,
		Name:        group.Name,
		Slug:        group.Slug,
		Description: sqlNullString(group.Description),
		Icon:        sqlNullString(group.Icon),
	})
	if err != nil {
		return nil, err
	}
	return sqlcToModelRecipeGroup(result), nil
}

func (r *RecipeGroupRepository) GetByID(id string) (*models.RecipeGroup, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	result, err := r.q.GetRecipeGroupByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	return sqlcToModelRecipeGroup(result), nil
}

func (r *RecipeGroupRepository) GetBySlug(slug string) (*models.RecipeGroup, error) {
	ctx := context.Background()
	result, err := r.q.GetRecipeGroupBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return sqlcToModelRecipeGroup(result), nil
}

func (r *RecipeGroupRepository) List() ([]*models.RecipeGroup, error) {
	ctx := context.Background()
	results, err := r.q.ListRecipeGroups(ctx)
	if err != nil {
		return nil, err
	}
	groups := make([]*models.RecipeGroup, len(results))
	for i, row := range results {
		groups[i] = sqlcToModelRecipeGroup(row)
	}
	return groups, nil
}

func (r *RecipeGroupRepository) Update(id string, group *models.RecipeGroup) (*models.RecipeGroup, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	result, err := r.q.UpdateRecipeGroup(ctx, sqlc.UpdateRecipeGroupParams{
		ID:          parsedID,
		Name:        sqlNullStringPtr(group.Name),
		Slug:        sqlNullStringPtr(group.Slug),
		Description: sqlNullString(group.Description),
		Icon:        sqlNullString(group.Icon),
	})
	if err != nil {
		return nil, err
	}
	return sqlcToModelRecipeGroup(result), nil
}

func (r *RecipeGroupRepository) Delete(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteRecipeGroup(ctx, parsedID)
}

func (r *RecipeGroupRepository) AddRecipeToGroup(groupID, recipeID string) error {
	ctx := context.Background()
	parsedGroupID, err := uuid.Parse(groupID)
	if err != nil {
		return err
	}
	parsedRecipeID, err := uuid.Parse(recipeID)
	if err != nil {
		return err
	}
	order := int32(0)
	orderNull := sqlNullInt32Ptr(order)
	return r.q.AddRecipeToGroup(ctx, sqlc.AddRecipeToGroupParams{
		GroupID:    parsedGroupID,
		RecipeID:   parsedRecipeID,
		OrderIndex: orderNull,
	})
}

func (r *RecipeGroupRepository) RemoveRecipeFromGroup(groupID, recipeID string) error {
	ctx := context.Background()
	parsedGroupID, err := uuid.Parse(groupID)
	if err != nil {
		return err
	}
	parsedRecipeID, err := uuid.Parse(recipeID)
	if err != nil {
		return err
	}
	return r.q.RemoveRecipeFromGroup(ctx, sqlc.RemoveRecipeFromGroupParams{
		GroupID:  parsedGroupID,
		RecipeID: parsedRecipeID,
	})
}

func (r *RecipeGroupRepository) GetRecipesInGroup(groupID string) ([]*models.Recipe, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(groupID)
	if err != nil {
		return nil, err
	}

	results, err := r.q.GetRecipesInGroup(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	recipes := make([]*models.Recipe, len(results))
	for i, row := range results {
		recipes[i] = &models.Recipe{
			ID:                row.ID,
			Title:             row.Title,
			Slug:              row.Slug,
			Description:       nullStringToPtr(row.Description),
			FeaturedImagePath: nullStringToPtr(row.FeaturedImagePath),
			MarkdownContent:   row.MarkdownContent,
			IsPublished:       row.IsPublished.Bool,
			CategoryID:        nullUUIDToPtr(row.CategoryID),
			AuthorID:          nullUUIDToPtr(row.AuthorID),
			CreatedAt:         row.CreatedAt.Time,
			UpdatedAt:         row.UpdatedAt.Time,
			PublishedAt:       nullTimeToTimePtr(row.PublishedAt),
		}
	}
	return recipes, nil
}

func sqlcToModelRecipeGroup(row sqlc.RecipeGroup) *models.RecipeGroup {
	id := row.ID
	return &models.RecipeGroup{
		ID:          id,
		Name:        row.Name,
		Slug:        row.Slug,
		Description: nullStringToPtr(row.Description),
		Icon:        nullStringToPtr(row.Icon),
		CreatedAt:   row.CreatedAt.Time,
	}
}
