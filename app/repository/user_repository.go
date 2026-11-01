package repository

import (
	"context"
	"errors"
	"fmt"

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
			return nil, ErrNotFound
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
			return nil, ErrNotFound
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

func (r *UserRepository) FindUsersByCourseID(ctx context.Context, courseID int) ([]model.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.role, u.is_active, u.created_at
		FROM users u
		JOIN enrollments e ON e.user_id = u.id
		WHERE e.course_id = $1 AND u.is_active = true
		ORDER BY u.id ASC
	`
	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("mengambil peserta kelas: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("membaca row user: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) LoadPermissions(ctx context.Context) (*helper.PermissionSet, error) {
	rows, err := r.db.Query(ctx, "SELECT role_name, permission_name FROM role_permissions")
	if err != nil {
		return nil, fmt.Errorf("memuat permission: %w", err)
	}
	defer rows.Close()

	raw := make(map[string][]string)
	for rows.Next() {
		var role, perm string
		if err := rows.Scan(&role, &perm); err != nil {
			return nil, err
		}
		raw[role] = append(raw[role], perm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return helper.NewPermissionSet(raw), nil
}
