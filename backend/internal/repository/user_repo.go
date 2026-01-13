package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/homecooking/backend/internal/db/sqlc"
	"github.com/homecooking/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewUserRepository(db *sql.DB, q *sqlc.Queries) *UserRepository {
	return &UserRepository{
		db: db,
		q:  q,
	}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	ctx := context.Background()
	result, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return nil, err
	}

	dbUser, err := r.q.GetUserByID(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return r.sqlcToModel(dbUser), nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	ctx := context.Background()
	dbUser, err := r.q.GetUserByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(dbUser), nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(dbUser), nil
}

func (r *UserRepository) List(limit, offset int) ([]*models.User, error) {
	ctx := context.Background()
	dbUsers, err := r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.sqlcToModel(dbUser)
	}
	return users, nil
}

func (r *UserRepository) Update(id string, user *models.User) (*models.User, error) {
	ctx := context.Background()
	result, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    uuid.MustParse(id),
		Email: sqlString(user.Email),
		Role:  sqlString(user.Role),
	})
	if err != nil {
		return nil, err
	}
	return r.sqlcToModel(result), nil
}

func (r *UserRepository) Delete(id string) error {
	ctx := context.Background()
	return r.q.DeleteUser(ctx, uuid.MustParse(id))
}

func (r *UserRepository) sqlcToModel(dbUser sqlc.User) *models.User {
	return &models.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}
}
