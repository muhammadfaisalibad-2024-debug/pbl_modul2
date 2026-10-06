package repository

import (
	"context"
	"errors"
	"time"

	"api-students/app/model"
	"api-students/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func (r *TokenRepository) Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`,
		userID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return helper.Internal(err)
	}
	return nil
}

func (r *TokenRepository) FindValid(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.pool.QueryRow(
		ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`,
		tokenHash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.Unauthorized("refresh token tidak valid")
		}
		return nil, helper.Internal(err)
	}
	return &token, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`,
		id,
	)
	if err != nil {
		return helper.Internal(err)
	}
	return nil
}

func (r *TokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`,
		tokenHash,
	)
	if err != nil {
		return helper.Internal(err)
	}
	return nil
}
