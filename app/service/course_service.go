package service

import (
	"strconv"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type CourseService struct {
	repo *repository.CourseRepository
}

func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{
		repo: repo,
	}
}

// Endpoint 8: GET /api/v1/courses (Semua role terautentikasi)
func (s *CourseService) GetCourses(c *fiber.Ctx) error {
	semester, _ := strconv.Atoi(c.Query("semester", "0"))
	search := strings.TrimSpace(c.Query("search"))
	available := strings.ToLower(c.Query("available")) == "true"

	query := model.CourseQuery{
		Semester:  semester,
		Search:    search,
		Available: available,
	}

	courses, err := s.repo.List(c.Context(), query)
	if err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusOK, "Daftar mata kuliah berhasil diambil", courses)
}
