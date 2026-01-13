package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type ShareCodeRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewShareCodeRepository(db *sql.DB, q *sqlc.Queries) *ShareCodeRepository {
	return &ShareCodeRepository{
		db: db,
		q:  q,
	}
}

func (r *ShareCodeRepository) Create(shareCode *models.ShareCode) (*models.ShareCode, error) {
	ctx := context.Background()

	var maxUses sql.NullInt32
	if shareCode.MaxUses != nil {
		val := int32(*shareCode.MaxUses)
		maxUses = sql.NullInt32{Int32: val, Valid: true}
	}

	result, err := r.q.CreateShareCode(ctx, sqlc.CreateShareCodeParams{
		RecipeID:  uuid.NullUUID{UUID: shareCode.RecipeID, Valid: true},
		Code:      shareCode.Code,
		ExpiresAt: sqlNullTimePtr(shareCode.ExpiresAt),
		MaxUses:   maxUses,
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *ShareCodeRepository) GetByCode(code string) (*models.ShareCodeWithRecipe, error) {
	ctx := context.Background()
	result, err := r.q.GetShareCodeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModelWithRecipe(result), nil
}

func (r *ShareCodeRepository) GetForRecipe(recipeID string) ([]*models.ShareCode, error) {
	ctx := context.Background()
	results, err := r.q.GetShareCodesForRecipe(ctx, uuid.NullUUID{UUID: uuid.MustParse(recipeID), Valid: true})
	if err != nil {
		return nil, err
	}

	shareCodes := make([]*models.ShareCode, len(results))
	for i, result := range results {
		shareCodes[i] = r.sqlcToModel(result)
	}
	return shareCodes, nil
}

func (r *ShareCodeRepository) IncrementUse(id string) error {
	ctx := context.Background()
	return r.q.IncrementShareCodeUse(ctx, uuid.MustParse(id))
}

func (r *ShareCodeRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteShareCode(ctx, uuid.MustParse(id))
}

func (r *ShareCodeRepository) sqlcToModel(dbShareCode sqlc.ShareCode) *models.ShareCode {
	useCount := int(dbShareCode.UseCount.Int32)
	maxUses := int32PtrToInt(dbShareCode.MaxUses.Int32)
	return &models.ShareCode{
		ID:        dbShareCode.ID,
		RecipeID:  dbShareCode.RecipeID.UUID,
		Code:      dbShareCode.Code,
		ExpiresAt: nullTimeToTimePtr(dbShareCode.ExpiresAt),
		MaxUses:   maxUses,
		UseCount:  useCount,
		CreatedAt: dbShareCode.CreatedAt.Time,
	}
}

func (r *ShareCodeRepository) sqlcToModelWithRecipe(dbRow sqlc.GetShareCodeByCodeRow) *models.ShareCodeWithRecipe {
	useCount := int(dbRow.UseCount.Int32)
	maxUses := int32PtrToInt(dbRow.MaxUses.Int32)
	return &models.ShareCodeWithRecipe{
		ID:          dbRow.ID,
		RecipeID:    dbRow.RecipeID.UUID,
		Code:        dbRow.Code,
		ExpiresAt:   nullTimeToTimePtr(dbRow.ExpiresAt),
		MaxUses:     maxUses,
		UseCount:    useCount,
		CreatedAt:   dbRow.CreatedAt.Time,
		RecipeTitle: dbRow.RecipeTitle,
		RecipeSlug:  dbRow.RecipeSlug,
	}
}
