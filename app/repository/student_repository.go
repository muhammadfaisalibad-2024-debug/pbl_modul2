package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"api-students/app/model"
	"api-students/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) List(ctx context.Context, q model.StudentListQuery) ([]model.StudentListItem, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	conditions = append(conditions, "deleted_at IS NULL")

	if q.Prodi != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(prodi) = LOWER($%d)", argIdx))
		args = append(args, q.Prodi)
		argIdx++
	}

	if q.Angkatan > 0 {
		conditions = append(conditions, fmt.Sprintf("angkatan = $%d", argIdx))
		args = append(args, q.Angkatan)
		argIdx++
	}

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(nim ILIKE '%%' || $%d || '%%' OR nama ILIKE '%%' || $%d || '%%')", argIdx, argIdx))
		args = append(args, q.Search)
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, helper.Internal(err)
	}

	// Order by
	orderBy := "nama ASC, id ASC"
	switch q.Sort {
	case "-ipk_terakhir":
		orderBy = "ipk_terakhir DESC, id ASC"
	case "ipk_terakhir":
		orderBy = "ipk_terakhir ASC, id ASC"
	case "-nama":
		orderBy = "nama DESC, id ASC"
	case "nama":
		orderBy = "nama ASC, id ASC"
	}

	if q.PerPage <= 0 {
		q.PerPage = 10
	}
	if q.PerPage > 50 {
		q.PerPage = 50
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	offset := (q.Page - 1) * q.PerPage

	selectQuery := fmt.Sprintf(`
		SELECT id, nim, nama, prodi, angkatan, COALESCE(ipk_terakhir, 0.00)
		FROM students
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIdx, argIdx+1)

	args = append(args, q.PerPage, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, helper.Internal(err)
	}
	defer rows.Close()

	items := make([]model.StudentListItem, 0)
	for rows.Next() {
		var item model.StudentListItem
		err := rows.Scan(
			&item.ID,
			&item.NIM,
			&item.Nama,
			&item.Prodi,
			&item.Angkatan,
			&item.IPKTerakhir,
		)
		if err != nil {
			return nil, 0, helper.Internal(err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, helper.Internal(err)
	}

	return items, total, nil
}

func (r *StudentRepository) FindByID(ctx context.Context, id int) (*model.Student, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), nim, nama, prodi, angkatan, COALESCE(ipk_terakhir, 0.00), deleted_at, created_at
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s model.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.NIM,
		&s.Nama,
		&s.Prodi,
		&s.Angkatan,
		&s.IPKTerakhir,
		&s.DeletedAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("data mahasiswa tidak ditemukan atau sudah dihapus")
		}
		return nil, helper.Internal(err)
	}
	return &s, nil
}

func (r *StudentRepository) FindByUserID(ctx context.Context, userID int) (*model.Student, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), nim, nama, prodi, angkatan, COALESCE(ipk_terakhir, 0.00), deleted_at, created_at
		FROM students
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	var s model.Student
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.NIM,
		&s.Nama,
		&s.Prodi,
		&s.Angkatan,
		&s.IPKTerakhir,
		&s.DeletedAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &s, nil
}

func (r *StudentRepository) FindByNIM(ctx context.Context, nim string) (*model.Student, error) {
	query := `
		SELECT id, COALESCE(user_id, 0), nim, nama, prodi, angkatan, COALESCE(ipk_terakhir, 0.00), deleted_at, created_at
		FROM students
		WHERE nim = $1 AND deleted_at IS NULL
	`
	var s model.Student
	err := r.db.QueryRow(ctx, query, nim).Scan(
		&s.ID,
		&s.UserID,
		&s.NIM,
		&s.Nama,
		&s.Prodi,
		&s.Angkatan,
		&s.IPKTerakhir,
		&s.DeletedAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("data mahasiswa tidak ditemukan")
		}
		return nil, helper.Internal(err)
	}
	return &s, nil
}

func (r *StudentRepository) CreateWithUserInTx(ctx context.Context, req model.CreateStudentRequest, defaultPasswordHash string) (*model.Student, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, helper.Internal(err)
	}
	defer tx.Rollback(ctx)

	// 1. Check duplicate email in users
	var existingUserID int
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE LOWER(email) = LOWER($1)", req.Email).Scan(&existingUserID)
	if err == nil {
		return nil, helper.ValidationField("email", "Email sudah terdaftar")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, helper.Internal(err)
	}

	// 2. Check duplicate NIM in students
	var existingStudentID int
	err = tx.QueryRow(ctx, "SELECT id FROM students WHERE nim = $1", req.NIM).Scan(&existingStudentID)
	if err == nil {
		return nil, helper.ValidationField("nim", "NIM sudah terdaftar")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, helper.Internal(err)
	}

	// 3. Create user record (role mahasiswa, default password = hashed NIM)
	var userID int
	err = tx.QueryRow(ctx, `
		INSERT INTO users (username, email, password, role, is_active)
		VALUES ($1, $2, $3, 'mahasiswa', true)
		RETURNING id
	`, req.NIM, req.Email, defaultPasswordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") {
				return nil, helper.ValidationField("email", "Email sudah terdaftar")
			}
			return nil, helper.ValidationField("nim", "NIM sudah terdaftar")
		}
		return nil, helper.Internal(err)
	}

	// 4. Create student record
	ipk := 0.00
	if req.IPKTerakhir != nil {
		ipk = *req.IPKTerakhir
	}

	var s model.Student
	err = tx.QueryRow(ctx, `
		INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, created_at
	`, userID, req.NIM, req.Nama, req.Prodi, req.Angkatan, ipk).Scan(
		&s.ID,
		&s.UserID,
		&s.NIM,
		&s.Nama,
		&s.Prodi,
		&s.Angkatan,
		&s.IPKTerakhir,
		&s.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, helper.ValidationField("nim", "NIM sudah terdaftar")
		}
		return nil, helper.Internal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, helper.Internal(err)
	}

	return &s, nil
}

func (r *StudentRepository) Update(ctx context.Context, id int, req model.UpdateStudentRequest) (*model.Student, error) {
	ipkClause := "ipk_terakhir"
	args := []interface{}{req.Nama, req.Prodi, req.Angkatan, id}
	if req.IPKTerakhir != nil {
		ipkClause = "$5"
		args = []interface{}{req.Nama, req.Prodi, req.Angkatan, id, *req.IPKTerakhir}
	}

	query := fmt.Sprintf(`
		UPDATE students
		SET nama = $1, prodi = $2, angkatan = $3, ipk_terakhir = %s
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, COALESCE(user_id, 0), nim, nama, prodi, angkatan, COALESCE(ipk_terakhir, 0.00), deleted_at, created_at
	`, ipkClause)

	var s model.Student
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&s.ID,
		&s.UserID,
		&s.NIM,
		&s.Nama,
		&s.Prodi,
		&s.Angkatan,
		&s.IPKTerakhir,
		&s.DeletedAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helper.NotFound("data mahasiswa tidak ditemukan atau sudah dihapus")
		}
		return nil, helper.Internal(err)
	}
	return &s, nil
}

func (r *StudentRepository) SoftDelete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, `
		UPDATE students
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return helper.Internal(err)
	}
	if result.RowsAffected() == 0 {
		return helper.NotFound("data mahasiswa tidak ditemukan atau sudah dihapus")
	}
	return nil
}

func (r *StudentRepository) GetEnrolledCourses(ctx context.Context, studentID int) ([]model.EnrolledCourseItem, int, error) {
	query := `
		SELECT e.id, c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, e.tahun_akademik
		FROM enrollments e
		JOIN courses c ON c.id = e.course_id
		WHERE e.student_id = $1
		ORDER BY c.semester ASC, c.kode_mk ASC
	`
	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, 0, helper.Internal(err)
	}
	defer rows.Close()

	items := make([]model.EnrolledCourseItem, 0)
	totalSKS := 0
	for rows.Next() {
		var item model.EnrolledCourseItem
		err := rows.Scan(
			&item.EnrollmentID,
			&item.CourseID,
			&item.KodeMK,
			&item.NamaMK,
			&item.SKS,
			&item.Semester,
			&item.TahunAkademik,
		)
		if err != nil {
			return nil, 0, helper.Internal(err)
		}
		totalSKS += item.SKS
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, helper.Internal(err)
	}
	return items, totalSKS, nil
}
