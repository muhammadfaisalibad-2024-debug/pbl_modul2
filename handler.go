package main

import (
	"strconv"

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

// GET /api/v1/students
func getStudents(c *fiber.Ctx) error {
	return sendSuccess(c, fiber.StatusOK, "Students retrieved successfully", students)
}

// GET /api/v1/students/:id
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

// POST /api/v1/students
func createStudent(c *fiber.Ctx) error {
	var request CreateStudentRequest

	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if request.NIM == "" || request.Name == "" {
		return sendError(c, fiber.StatusUnprocessableEntity, "NIM and Name are required")
	}

	// Cek NIM duplikat
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
