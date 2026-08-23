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

func getStudents(c *fiber.Ctx) error {
	return sendSuccess(c, fiber.StatusOK, "Students retrieved successfully", students)
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
