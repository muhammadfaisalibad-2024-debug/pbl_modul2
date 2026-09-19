package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"api-students/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	FindByUsername(ctx context.Context, username string) (model.User, error)
	Update(ctx context.Context, id int, req model.ReplaceUserRequest) (model.User, error)
	Patch(ctx context.Context, id int, req model.PatchUserRequest) (model.User, error)
	UpdateRole(ctx context.Context, id int, role string) (model.User, error)
	Delete(ctx context.Context, id int) error
}

func (r *userPostgresRepository) FindAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, username, email, password, role, is_active, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar user: %w", err)
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("membaca user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca daftar user: %w", err)
	}
	return users, nil
}

type userPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userPostgresRepository{pool: pool}
}

func (r *userPostgresRepository) FindByUsername(
	ctx context.Context,
	username string,
) (model.User, error) {

	var u model.User

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, username, email, password, role, is_active, created_at
		FROM users
		WHERE LOWER(username) = LOWER($1)`,
		username,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, fmt.Errorf("mengambil user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) FindByID(
	ctx context.Context,
	id int,
) (model.User, error) {

	var u model.User

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, username, email, password, role, is_active, created_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, fmt.Errorf("mengambil user: %w", err)
	}

	return u, nil
}

func (r *userPostgresRepository) Create(
	ctx context.Context,
	user model.User,
) (model.User, error) {

	var created model.User

	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO users
		(username, email, password, role, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, password, role, is_active, created_at`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
	).Scan(
		&created.ID,
		&created.Username,
		&created.Email,
		&created.Password,
		&created.Role,
		&created.IsActive,
		&created.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return model.User{}, ErrDuplicate
		}
		return model.User{}, err
	}

	return created, nil
}

func (r *userPostgresRepository) Update(ctx context.Context, id int, req model.ReplaceUserRequest) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `
		UPDATE users SET username = $1, email = $2, is_active = $3
		WHERE id = $4
		RETURNING id, username, email, password, role, is_active, created_at`,
		req.Username, req.Email, req.IsActive, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	return userMutationResult(u, err, "mengubah user")
}

func (r *userPostgresRepository) Patch(ctx context.Context, id int, req model.PatchUserRequest) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `
		UPDATE users SET
			username = COALESCE($1, username),
			email = COALESCE($2, email),
			is_active = COALESCE($3, is_active)
		WHERE id = $4
		RETURNING id, username, email, password, role, is_active, created_at`,
		req.Username, req.Email, req.IsActive, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	return userMutationResult(u, err, "mengubah user")
}

func (r *userPostgresRepository) UpdateRole(ctx context.Context, id int, role string) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `
		UPDATE users SET role = $1 WHERE id = $2
		RETURNING id, username, email, password, role, is_active, created_at`,
		role, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	return userMutationResult(u, err, "mengubah role user")
}

func (r *userPostgresRepository) Delete(ctx context.Context, id int) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("menghapus user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func userMutationResult(user model.User, err error, action string) (model.User, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", action, err)
	}
	return user, nil
}
