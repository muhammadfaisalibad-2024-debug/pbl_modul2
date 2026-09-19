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

// Sentinel errors for service layer mapping.
var (
	ErrNotFound  = errors.New("student not found")
	ErrDuplicate = errors.New("NIM already exists")
)

type StudentRepository interface {
	FindAll(ctx context.Context, page, limit int, search, sortBy, order, active string) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, req *model.CreateStudentRequest, ownerID int) (model.Student, error)
	Update(ctx context.Context, id int, req *model.ReplaceStudentRequest) (model.Student, error)
	Patch(ctx context.Context, id int, req *model.PatchStudentRequest) (model.Student, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type studentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindAll(ctx context.Context, page, limit int, search, sortBy, order, active string) ([]model.Student, int, error) {
	sortColumns := map[string]string{
		"id":    "id",
		"name":  "name",
		"grade": "grade",
	}

	sortColumn, ok := sortColumns[sortBy]
	if !ok {
		sortColumn = "id"
	}

	if order != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit

	where := []string{"1=1"}
	args := []interface{}{}
	argNumber := 1

	if search != "" {
		where = append(
			where,
			fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(nim) LIKE $%d)", argNumber, argNumber),
		)
		args = append(args, "%"+search+"%")
		argNumber++
	}

	if active != "" {
		isActive := active == "true"
		if active == "true" || active == "false" {
			where = append(where, fmt.Sprintf("is_active = $%d", argNumber))
			args = append(args, isActive)
			argNumber++
		}
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students WHERE %s", whereClause)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, nim, name, grade, is_active, COALESCE(owner_id, 0)
		FROM students
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortColumn, order, argNumber, argNumber+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := make([]model.Student, 0)
	for rows.Next() {
		var student model.Student
		if err := rows.Scan(
			&student.ID,
			&student.NIM,
			&student.Name,
			&student.Grade,
			&student.IsActive,
			&student.OwnerID,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, student)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *studentRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var student model.Student
	err := r.db.QueryRow(
		ctx,
		`SELECT id, nim, name, grade, is_active, COALESCE(owner_id, 0)
		 FROM students
		 WHERE id = $1`,
		id,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
		&student.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) Create(ctx context.Context, req *model.CreateStudentRequest, ownerID int) (model.Student, error) {
	var student model.Student
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO students (nim, name, grade, is_active, owner_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, nim, name, grade, is_active, owner_id`,
		req.NIM,
		req.Name,
		req.Grade,
		req.IsActive,
		ownerID,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
		&student.OwnerID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "students_nim_key") {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) Update(ctx context.Context, id int, req *model.ReplaceStudentRequest) (model.Student, error) {
	var student model.Student
	err := r.db.QueryRow(
		ctx,
		`UPDATE students
		 SET nim = $1,
		     name = $2,
		     grade = $3,
		     is_active = $4
		 WHERE id = $5
		 RETURNING id, nim, name, grade, is_active, COALESCE(owner_id, 0)`,
		req.NIM,
		req.Name,
		req.Grade,
		req.IsActive,
		id,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
		&student.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if strings.Contains(err.Error(), "students_nim_key") {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) Patch(ctx context.Context, id int, req *model.PatchStudentRequest) (model.Student, error) {
	setParts := []string{}
	args := []interface{}{}
	argNumber := 1

	if req.NIM != nil {
		setParts = append(setParts, fmt.Sprintf("nim = $%d", argNumber))
		args = append(args, *req.NIM)
		argNumber++
	}

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argNumber))
		args = append(args, *req.Name)
		argNumber++
	}

	if req.Grade != nil {
		setParts = append(setParts, fmt.Sprintf("grade = $%d", argNumber))
		args = append(args, *req.Grade)
		argNumber++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argNumber))
		args = append(args, *req.IsActive)
		argNumber++
	}

	if len(setParts) == 0 {
		return model.Student{}, fmt.Errorf("no fields to update")
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE students
		SET %s
		WHERE id = $%d
		RETURNING id, nim, name, grade, is_active, COALESCE(owner_id, 0)
	`, strings.Join(setParts, ", "), argNumber)

	var student model.Student
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
		&student.OwnerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if strings.Contains(err.Error(), "students_nim_key") {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(
		ctx,
		"DELETE FROM students WHERE id = $1",
		id,
	)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}
