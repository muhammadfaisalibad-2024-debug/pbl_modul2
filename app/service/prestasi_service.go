package service

import (
	"errors"

	"api-students/app/repository"
	"api-students/helper"
	"github.com/gofiber/fiber/v2"
)

type PrestasiService struct {
	students repository.StudentRepository
	prestasi repository.PrestasiRepository
}

func NewPrestasiService(students repository.StudentRepository, prestasi repository.PrestasiRepository) *PrestasiService {
	return &PrestasiService{students: students, prestasi: prestasi}
}

func (s *PrestasiService) GetByStudentNIM(c *fiber.Ctx) error {
	nim := c.Params("nim")
	student, err := s.students.FindByNIM(c.Context(), nim)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "Student not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to retrieve student")
	}

	prestasi, err := s.prestasi.FindByStudentID(c.Context(), student.ID)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to retrieve prestasi")
	}
	return helper.Success(c, fiber.StatusOK, "Prestasi retrieved successfully", prestasi)
}
