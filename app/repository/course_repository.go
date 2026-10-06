package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"api-students/app/model"
	"api-students/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) List(ctx context.Context, q model.CourseQuery) ([]model.CourseListItem, error) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if q.Semester > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("c.semester = $%d", argIdx))
		args = append(args, q.Semester)
		argIdx++
	}

	if q.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(c.kode_mk ILIKE '%%' || $%d || '%%' OR c.nama_mk ILIKE '%%' || $%d || '%%')", argIdx, argIdx))
		args = append(args, q.Search)
		argIdx++
	}

	whereSql := ""
	if len(whereClauses) > 0 {
		whereSql = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	havingSql := ""
	if q.Available {
		havingSql = "HAVING COUNT(e.id) < c.kuota"
	}

	query := fmt.Sprintf(`
		SELECT 
			c.id, 
			c.kode_mk, 
			c.nama_mk, 
			c.sks, 
			c.semester, 
			c.kuota,
			COUNT(e.id)::int AS terisi,
			(c.kuota - COUNT(e.id))::int AS sisa_kuota
		FROM courses c
		LEFT JOIN enrollments e ON e.course_id = c.id
		%s
		GROUP BY c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota
		%s
		ORDER BY c.semester ASC, c.kode_mk ASC
	`, whereSql, havingSql)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, helper.Internal(err)
	}
	defer rows.Close()

	items := make([]model.CourseListItem, 0)
	for rows.Next() {
		var item model.CourseListItem
		err := rows.Scan(
			&item.ID,
			&item.KodeMK,
			&item.NamaMK,
			&item.SKS,
			&item.Semester,
			&item.Kuota,
			&item.Terisi,
			&item.SisaKuota,
		)
		if err != nil {
			return nil, helper.Internal(err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, helper.Internal(err)
	}

	return items, nil
}

func (r *CourseRepository) FindByID(ctx context.Context, id int) (*model.Course, error) {
	query := `
		SELECT id, kode_mk, nama_mk, sks, semester, kuota, created_at
		FROM courses
		WHERE id = $1
	`
	var c model.Course
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.KodeMK,
		&c.NamaMK,
		&c.SKS,
		&c.Semester,
		&c.Kuota,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("mata kuliah tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &c, nil
}
