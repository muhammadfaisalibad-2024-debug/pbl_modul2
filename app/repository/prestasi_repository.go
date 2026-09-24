package repository

import (
	"context"
	"fmt"

	"api-students/app/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrestasiRepository interface {
	FindByStudentID(ctx context.Context, studentID int) ([]model.Prestasi, error)
}

type prestasiPostgresRepository struct{ db *pgxpool.Pool }

func NewPrestasiRepository(db *pgxpool.Pool) PrestasiRepository {
	return &prestasiPostgresRepository{db: db}
}

func (r *prestasiPostgresRepository) FindByStudentID(ctx context.Context, studentID int) ([]model.Prestasi, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, id_student, nama, juara
		FROM prestasi
		WHERE id_student = $1
		ORDER BY id`, studentID)
	if err != nil {
		return nil, fmt.Errorf("mengambil prestasi: %w", err)
	}
	defer rows.Close()

	result := make([]model.Prestasi, 0)
	for rows.Next() {
		var prestasi model.Prestasi
		if err := rows.Scan(&prestasi.ID, &prestasi.StudentID, &prestasi.Nama, &prestasi.Juara); err != nil {
			return nil, fmt.Errorf("membaca prestasi: %w", err)
		}
		result = append(result, prestasi)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil prestasi: %w", err)
	}
	return result, nil
}
