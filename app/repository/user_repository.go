package repository

import (
	"context"
	"errors"

	"api-students/app/model"
	"api-students/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, role, is_active, created_at
		FROM users
		WHERE LOWER(email) = LOWER($1) AND is_active = true
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("user tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, username, email, password, role, is_active, created_at
		FROM users
		WHERE id = $1 AND is_active = true
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("user tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, role, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Role).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return helper.Internal(err)
	}
	return nil
}
