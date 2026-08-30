package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func getStudents(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	search := strings.ToLower(c.Query("search"))
	sortBy := c.Query("sort", "id")
	order := strings.ToLower(c.Query("order", "asc"))
	active := c.Query("is_active")

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
			fmt.Sprintf(
				"(LOWER(name) LIKE $%d OR LOWER(nim) LIKE $%d)",
				argNumber,
				argNumber,
			),
		)
		args = append(args, "%"+search+"%")
		argNumber++
	}

	if active != "" {
		isActive, err := strconv.ParseBool(active)

		if err == nil {
			where = append(
				where,
				fmt.Sprintf("is_active = $%d", argNumber),
			)
			args = append(args, isActive)
			argNumber++
		}
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM students WHERE %s",
		whereClause,
	)

	var total int

	if err := db.QueryRow(c.Context(), countQuery, args...).Scan(&total); err != nil {
		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve students",
		)
	}

	query := fmt.Sprintf(`
		SELECT id, nim, name, grade, is_active
		FROM students
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortColumn, order, argNumber, argNumber+1)

	args = append(args, limit, offset)

	rows, err := db.Query(c.Context(), query, args...)
	if err != nil {
		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve students",
		)
	}
	defer rows.Close()

	result := make([]Student, 0)

	for rows.Next() {
		var student Student

		if err := rows.Scan(
			&student.ID,
			&student.NIM,
			&student.Name,
			&student.Grade,
			&student.IsActive,
		); err != nil {
			return sendError(
				c,
				fiber.StatusInternalServerError,
				"Failed to read student data",
			)
		}

		result = append(result, student)
	}

	if err := rows.Err(); err != nil {
		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to read student data",
		)
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Students retrieved successfully",
		"data":    result,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func getStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	var student Student

	err = db.QueryRow(
		c.Context(),
		`SELECT id, nim, name, grade, is_active
		 FROM students
		 WHERE id = $1`,
		id,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
	)

	if err == pgx.ErrNoRows {
		return sendError(c, fiber.StatusNotFound, "Student not found")
	}

	if err != nil {
		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to retrieve student",
		)
	}

	return sendSuccess(
		c,
		fiber.StatusOK,
		"Student retrieved successfully",
		student,
	)
}

func createStudent(c *fiber.Ctx) error {
	var request CreateStudentRequest

	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if request.NIM == "" || request.Name == "" {
		return sendError(
			c,
			fiber.StatusUnprocessableEntity,
			"NIM and Name are required",
		)
	}

	var student Student

	err := db.QueryRow(
		c.Context(),
		`INSERT INTO students (nim, name, grade, is_active)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, nim, name, grade, is_active`,
		request.NIM,
		request.Name,
		request.Grade,
		request.IsActive,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
	)

	if err != nil {
		if strings.Contains(err.Error(), "students_nim_key") {
			return sendError(c, fiber.StatusConflict, "NIM already exists")
		}

		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to create student",
		)
	}

	c.Location("/api/v1/students/" + strconv.Itoa(student.ID))

	return sendSuccess(
		c,
		fiber.StatusCreated,
		"Student created successfully",
		student,
	)
}

func updateStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	var request ReplaceStudentRequest

	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if request.NIM == "" || request.Name == "" {
		return sendError(
			c,
			fiber.StatusUnprocessableEntity,
			"NIM and Name are required",
		)
	}

	var student Student

	err = db.QueryRow(
		c.Context(),
		`UPDATE students
		 SET nim = $1,
		     name = $2,
		     grade = $3,
		     is_active = $4
		 WHERE id = $5
		 RETURNING id, nim, name, grade, is_active`,
		request.NIM,
		request.Name,
		request.Grade,
		request.IsActive,
		id,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
	)

	if err == pgx.ErrNoRows {
		return sendError(c, fiber.StatusNotFound, "Student not found")
	}

	if err != nil {
		if strings.Contains(err.Error(), "students_nim_key") {
			return sendError(c, fiber.StatusConflict, "NIM already exists")
		}

		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to update student",
		)
	}

	return sendSuccess(
		c,
		fiber.StatusOK,
		"Student updated successfully",
		student,
	)
}

func patchStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	var request PatchStudentRequest

	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	setParts := []string{}
	args := []interface{}{}
	argNumber := 1

	if request.NIM != nil {
		setParts = append(setParts, fmt.Sprintf("nim = $%d", argNumber))
		args = append(args, *request.NIM)
		argNumber++
	}

	if request.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argNumber))
		args = append(args, *request.Name)
		argNumber++
	}

	if request.Grade != nil {
		setParts = append(setParts, fmt.Sprintf("grade = $%d", argNumber))
		args = append(args, *request.Grade)
		argNumber++
	}

	if request.IsActive != nil {
		setParts = append(
			setParts,
			fmt.Sprintf("is_active = $%d", argNumber),
		)
		args = append(args, *request.IsActive)
		argNumber++
	}

	if len(setParts) == 0 {
		return sendError(
			c,
			fiber.StatusUnprocessableEntity,
			"No fields to update",
		)
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE students
		SET %s
		WHERE id = $%d
		RETURNING id, nim, name, grade, is_active
	`, strings.Join(setParts, ", "), argNumber)

	var student Student

	err = db.QueryRow(
		c.Context(),
		query,
		args...,
	).Scan(
		&student.ID,
		&student.NIM,
		&student.Name,
		&student.Grade,
		&student.IsActive,
	)

	if err == pgx.ErrNoRows {
		return sendError(c, fiber.StatusNotFound, "Student not found")
	}

	if err != nil {
		if strings.Contains(err.Error(), "students_nim_key") {
			return sendError(c, fiber.StatusConflict, "NIM already exists")
		}

		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to patch student",
		)
	}

	return sendSuccess(
		c,
		fiber.StatusOK,
		"Student patched successfully",
		student,
	)
}

func deleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	result, err := db.Exec(
		c.Context(),
		"DELETE FROM students WHERE id = $1",
		id,
	)

	if err != nil {
		return sendError(
			c,
			fiber.StatusInternalServerError,
			"Failed to delete student",
		)
	}

	if result.RowsAffected() == 0 {
		return sendError(c, fiber.StatusNotFound, "Student not found")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
