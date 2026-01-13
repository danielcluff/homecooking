package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type UserInviteRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewUserInviteRepository(db *sql.DB, q *sqlc.Queries) *UserInviteRepository {
	return &UserInviteRepository{
		db: db,
		q:  q,
	}
}

func (r *UserInviteRepository) Create(invite *models.UserInvite) (*models.UserInvite, error) {
	ctx := context.Background()

	result, err := r.q.CreateUserInvite(ctx, sqlc.CreateUserInviteParams{
		Code:      invite.Code,
		Email:     sqlNullString(invite.Email),
		Role:      sqlString(invite.Role),
		CreatedBy: sqlNullUUID(invite.CreatedBy),
		ExpiresAt: sqlNullTimePtr(invite.ExpiresAt),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *UserInviteRepository) GetByCode(code string) (*models.UserInvite, error) {
	ctx := context.Background()
	result, err := r.q.GetInviteByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *UserInviteRepository) List() ([]*models.UserInvite, error) {
	ctx := context.Background()
	results, err := r.q.ListInvites(ctx)
	if err != nil {
		return nil, err
	}

	invites := make([]*models.UserInvite, len(results))
	for i, result := range results {
		invites[i] = r.sqlcToModel(result)
	}
	return invites, nil
}

func (r *UserInviteRepository) Use(id string, usedBy string) (*models.UserInvite, error) {
	ctx := context.Background()
	usedByUUID := uuid.MustParse(usedBy)
	result, err := r.q.UseInvite(ctx, sqlc.UseInviteParams{
		ID:     uuid.MustParse(id),
		UsedBy: sqlNullUUID(&usedByUUID),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *UserInviteRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteInvite(ctx, uuid.MustParse(id))
}

func (r *UserInviteRepository) sqlcToModel(dbInvite sqlc.UserInvite) *models.UserInvite {
	role := "user"
	if dbInvite.Role.Valid {
		role = dbInvite.Role.String
	}

	return &models.UserInvite{
		ID:        dbInvite.ID,
		Code:      dbInvite.Code,
		Email:     nullStringToPtr(dbInvite.Email),
		Role:      role,
		CreatedBy: nullUUIDToPtr(dbInvite.CreatedBy),
		ExpiresAt: nullTimeToTimePtr(dbInvite.ExpiresAt),
		UsedAt:    nullTimeToTimePtr(dbInvite.UsedAt),
		UsedBy:    nullUUIDToPtr(dbInvite.UsedBy),
		CreatedAt: dbInvite.CreatedAt.Time,
	}
}
