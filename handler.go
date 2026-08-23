package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var students = []Student{
	{
		ID:       1,
		NIM:      "0123456789",
		Name:     "Faisal",
		Grade:    90,
		IsActive: true,
	},
	{
		ID:       2,
		NIM:      "0123456788",
		Name:     "Budi",
		Grade:    85,
		IsActive: true,
	},
}

var nextID = 3

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
	order := c.Query("order", "asc")
	active := c.Query("is_active")

	result := make([]Student, 0)

	for _, student := range students {
		if search != "" {
			nameMatch := strings.Contains(
				strings.ToLower(student.Name),
				search,
			)

			nimMatch := strings.Contains(
				strings.ToLower(student.NIM),
				search,
			)

			if !nameMatch && !nimMatch {
				continue
			}
		}

		if active != "" {
			isActive, err := strconv.ParseBool(active)

			if err == nil && student.IsActive != isActive {
				continue
			}
		}

		result = append(result, student)
	}

	sort.Slice(result, func(i, j int) bool {
		switch sortBy {
		case "name":
			if order == "desc" {
				return result[i].Name > result[j].Name
			}
			return result[i].Name < result[j].Name

		case "grade":
			if order == "desc" {
				return result[i].Grade > result[j].Grade
			}
			return result[i].Grade < result[j].Grade

		default:
			if order == "desc" {
				return result[i].ID > result[j].ID
			}
			return result[i].ID < result[j].ID
		}
	})

	total := len(result)

	start := (page - 1) * limit
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	paginated := result[start:end]

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Students retrieved successfully",
		"data":    paginated,
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

	for _, student := range students {
		if student.ID == id {
			return sendSuccess(c, fiber.StatusOK, "Student retrieved successfully", student)
		}
	}

	return sendError(c, fiber.StatusNotFound, "Student not found")
}

func createStudent(c *fiber.Ctx) error {
	var request CreateStudentRequest

	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if request.NIM == "" || request.Name == "" {
		return sendError(c, fiber.StatusUnprocessableEntity, "NIM and Name are required")
	}

	for _, student := range students {
		if student.NIM == request.NIM {
			return sendError(c, fiber.StatusConflict, "NIM already exists")
		}
	}

	student := Student{
		ID:       nextID,
		NIM:      request.NIM,
		Name:     request.Name,
		Grade:    request.Grade,
		IsActive: request.IsActive,
	}

	students = append(students, student)
	nextID++

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

	for i := range students {
		if students[i].ID == id {

			for _, student := range students {
				if student.ID != id && student.NIM == request.NIM {
					return sendError(
						c,
						fiber.StatusConflict,
						"NIM already exists",
					)
				}
			}

			students[i].NIM = request.NIM
			students[i].Name = request.Name
			students[i].Grade = request.Grade
			students[i].IsActive = request.IsActive

			return sendSuccess(
				c,
				fiber.StatusOK,
				"Student updated successfully",
				students[i],
			)
		}
	}

	return sendError(c, fiber.StatusNotFound, "Student not found")
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

	for i := range students {
		if students[i].ID == id {

			if request.NIM != nil {
				for _, student := range students {
					if student.ID != id && student.NIM == *request.NIM {
						return sendError(
							c,
							fiber.StatusConflict,
							"NIM already exists",
						)
					}
				}

				students[i].NIM = *request.NIM
			}

			if request.Name != nil {
				students[i].Name = *request.Name
			}

			if request.Grade != nil {
				students[i].Grade = *request.Grade
			}

			if request.IsActive != nil {
				students[i].IsActive = *request.IsActive
			}

			return sendSuccess(
				c,
				fiber.StatusOK,
				"Student patched successfully",
				students[i],
			)
		}
	}

	return sendError(c, fiber.StatusNotFound, "Student not found")
}

func deleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	for i, student := range students {
		if student.ID == id {

			students = append(
				students[:i],
				students[i+1:]...,
			)

			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return sendError(c, fiber.StatusNotFound, "Student not found")
}
