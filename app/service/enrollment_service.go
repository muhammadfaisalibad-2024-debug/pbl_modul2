package service

import (
	"strconv"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type EnrollmentService struct {
	enrollmentRepo *repository.EnrollmentRepository
	studentRepo    *repository.StudentRepository
}

func NewEnrollmentService(
	enrollmentRepo *repository.EnrollmentRepository,
	studentRepo *repository.StudentRepository,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		studentRepo:    studentRepo,
	}
}

// Endpoint 9: POST /api/v1/enrollments (Mahasiswa only)
func (s *EnrollmentService) CreateEnrollment(c *fiber.Ctx) error {
	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	if authUser.Role != model.RoleMahasiswa {
		return helper.Forbidden("hanya mahasiswa yang dapat mengambil mata kuliah (KRS)")
	}

	student, err := s.studentRepo.FindByUserID(c.Context(), authUser.UserID)
	if err != nil {
		return helper.Forbidden("data profil mahasiswa tidak ditemukan")
	}

	var req model.CreateEnrollmentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}

	if valErrors := helper.ValidateStruct(req); len(valErrors) > 0 {
		return helper.Validation(valErrors)
	}

	batasSKS := model.CalculateBatasSKS(student.IPKTerakhir)

	enrollment, err := s.enrollmentRepo.CreateEnrollment(c.Context(), student.ID, req.CourseID, req.TahunAkademik, batasSKS)
	if err != nil {
		return err
	}

	return helper.Created(c, "Mata kuliah berhasil ditambahkan ke KRS", enrollment)
}

// Endpoint 10: DELETE /api/v1/enrollments/:id (Mahasiswa milik sendiri)
func (s *EnrollmentService) DeleteEnrollment(c *fiber.Ctx) error {
	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	if authUser.Role != model.RoleMahasiswa {
		return helper.Forbidden("hanya mahasiswa yang dapat membatalkan KRS")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.BadRequest("ID enrollment tidak valid")
	}

	student, err := s.studentRepo.FindByUserID(c.Context(), authUser.UserID)
	if err != nil {
		return helper.Forbidden("data profil mahasiswa tidak ditemukan")
	}

	if err := s.enrollmentRepo.DeleteEnrollment(c.Context(), id, student.ID); err != nil {
		return err
	}

	return helper.NoContent(c)
}
