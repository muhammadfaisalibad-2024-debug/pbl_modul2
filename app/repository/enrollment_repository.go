package repository

import (
	"context"
	"errors"
	"fmt"

	"api-students/app/model"
	"api-students/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) CreateEnrollment(ctx context.Context, studentID int, courseID int, tahunAkademik string, batasSKS int) (*model.Enrollment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, helper.Internal(err)
	}
	defer tx.Rollback(ctx)

	// 1. Lock course row and fetch course details (with FOR UPDATE)
	var course model.Course
	err = tx.QueryRow(ctx, `
		SELECT id, kode_mk, nama_mk, sks, semester, kuota
		FROM courses
		WHERE id = $1
		FOR UPDATE
	`, courseID).Scan(
		&course.ID,
		&course.KodeMK,
		&course.NamaMK,
		&course.SKS,
		&course.Semester,
		&course.Kuota,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.ValidationField("course_id", "Mata kuliah tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}

	// 2. Check duplicate enrollment in the same academic year
	var existingEnrollmentID int
	err = tx.QueryRow(ctx, `
		SELECT id FROM enrollments
		WHERE student_id = $1 AND course_id = $2 AND tahun_akademik = $3
	`, studentID, courseID, tahunAkademik).Scan(&existingEnrollmentID)
	if err == nil {
		return nil, helper.Conflict("Mata kuliah sudah pernah diambil pada tahun akademik ini")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, helper.Internal(err)
	}

	// 3. Check course quota
	var terisi int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM enrollments WHERE course_id = $1
	`, courseID).Scan(&terisi)
	if err != nil {
		return nil, helper.Internal(err)
	}

	if terisi >= course.Kuota {
		return nil, helper.ValidationField("kuota", "Kuota mata kuliah sudah penuh")
	}

	// 4. Check SKS limit for this student in this academic year
	var currentSKS int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(SUM(c.sks), 0)
		FROM enrollments e
		JOIN courses c ON c.id = e.course_id
		WHERE e.student_id = $1 AND e.tahun_akademik = $2
	`, studentID, tahunAkademik).Scan(&currentSKS)
	if err != nil {
		return nil, helper.Internal(err)
	}

	sisaSKS := batasSKS - currentSKS
	if sisaSKS < 0 {
		sisaSKS = 0
	}

	if currentSKS+course.SKS > batasSKS {
		msg := fmt.Sprintf("Total SKS melebihi batas maksimal. Sisa SKS yang dapat diambil: %d SKS, dibutuhkan: %d SKS (batas maksimal: %d SKS)", sisaSKS, course.SKS, batasSKS)
		return nil, helper.ValidationField("sks", msg)
	}

	// 5. Insert enrollment
	var enr model.Enrollment
	err = tx.QueryRow(ctx, `
		INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		VALUES ($1, $2, $3)
		RETURNING id, student_id, course_id, tahun_akademik, created_at
	`, studentID, courseID, tahunAkademik).Scan(
		&enr.ID,
		&enr.StudentID,
		&enr.CourseID,
		&enr.TahunAkademik,
		&enr.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, helper.Conflict("Mata kuliah sudah pernah diambil pada tahun akademik ini")
		}
		return nil, helper.Internal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, helper.Internal(err)
	}

	return &enr, nil
}

func (r *EnrollmentRepository) FindByID(ctx context.Context, id int) (*model.Enrollment, error) {
	query := `
		SELECT id, student_id, course_id, tahun_akademik, created_at
		FROM enrollments
		WHERE id = $1
	`
	var enr model.Enrollment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&enr.ID,
		&enr.StudentID,
		&enr.CourseID,
		&enr.TahunAkademik,
		&enr.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("data enrollment tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &enr, nil
}

func (r *EnrollmentRepository) DeleteEnrollment(ctx context.Context, id int, studentID int) error {
	// First check if enrollment exists
	var ownerStudentID int
	err := r.db.QueryRow(ctx, "SELECT student_id FROM enrollments WHERE id = $1", id).Scan(&ownerStudentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return helper.NotFound("data enrollment tidak ditemukan")
		}
		return helper.Internal(err)
	}

	if ownerStudentID != studentID {
		return helper.Forbidden("tidak memiliki akses untuk membatalkan KRS mahasiswa lain")
	}

	_, err = r.db.Exec(ctx, "DELETE FROM enrollments WHERE id = $1", id)
	if err != nil {
		return helper.Internal(err)
	}

	return nil
}
