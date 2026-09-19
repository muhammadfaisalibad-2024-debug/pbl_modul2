package service

import (
	"errors"
	"strconv"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	repo  repository.StudentRepository
	perms *helper.PermissionSet
}

func NewStudentService(repo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{repo: repo, perms: perms}
}

func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	q := helper.ParseListQuery(c)

	students, total, err := s.repo.FindAll(
		c.Context(),
		q.Page, q.Limit,
		q.Search, q.SortBy, q.Order, q.IsActive,
	)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to retrieve students")
	}

	meta := model.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: CountTotalPages(total, q.Limit),
	}

	return helper.SuccessList(c, "Students retrieved successfully", students, meta)
}

func (s *StudentService) GetStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "Student not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to retrieve student")
	}
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}
	if !CanAccessStudent(current, student.OwnerID, s.perms, "student:read:any") {
		return helper.Fail(c, fiber.StatusForbidden, "tidak berhak mengakses data student ini")
	}

	return helper.Success(c, fiber.StatusOK, "Student retrieved successfully", student)
}

func (s *StudentService) CreateStudent(c *fiber.Ctx) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}
	var req model.CreateStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if err := ValidateCreate(&req); err != nil {
		return helper.Fail(c, fiber.StatusUnprocessableEntity, "NIM and Name are required")
	}

	student, err := s.repo.Create(c.Context(), &req, current.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, "NIM already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to create student")
	}

	c.Location("/api/v1/students/" + strconv.Itoa(student.ID))
	return helper.Created(c, "Student created successfully", student)
}

func (s *StudentService) UpdateStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid student ID")
	}
	if err := s.authorizeStudentUpdate(c, id); err != nil {
		return err
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if err := ValidateReplace(&req); err != nil {
		return helper.Fail(c, fiber.StatusUnprocessableEntity, "NIM and Name are required")
	}

	student, err := s.repo.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "Student not found")
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, "NIM already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to update student")
	}

	return helper.Success(c, fiber.StatusOK, "Student updated successfully", student)
}

func (s *StudentService) PatchStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid student ID")
	}
	if err := s.authorizeStudentUpdate(c, id); err != nil {
		return err
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid JSON body")
	}

	if IsEmptyPatch(&req) {
		return helper.Fail(c, fiber.StatusUnprocessableEntity, "No fields to update")
	}

	student, err := s.repo.Patch(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "Student not found")
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, "NIM already exists")
		}
		// "no fields to update" guard from repository (fallback)
		if strings.Contains(err.Error(), "no fields to update") {
			return helper.Fail(c, fiber.StatusUnprocessableEntity, "No fields to update")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to patch student")
	}

	return helper.Success(c, fiber.StatusOK, "Student patched successfully", student)
}

func (s *StudentService) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "Invalid student ID")
	}

	deleted, err := s.repo.Delete(c.Context(), id)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to delete student")
	}

	if !deleted {
		return helper.Fail(c, fiber.StatusNotFound, "Student not found")
	}

	return helper.NoContent(c)
}

func (s *StudentService) authorizeStudentUpdate(c *fiber.Ctx, id int) error {
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}
	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "Student not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "Failed to retrieve student")
	}
	if !CanAccessStudent(current, student.OwnerID, s.perms, "student:update:any") {
		return helper.Fail(c, fiber.StatusForbidden, "tidak berhak mengubah data student ini")
	}
	return nil
}
