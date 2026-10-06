package service

import (
	"strconv"
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

// Endpoint 3: GET /api/v1/students (Admin only)
func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage <= 0 {
		perPage = 10
	}
	if perPage > 50 {
		perPage = 50
	}

	angkatan, _ := strconv.Atoi(c.Query("angkatan", "0"))

	query := model.StudentListQuery{
		Page:     page,
		PerPage:  perPage,
		Prodi:    c.Query("prodi"),
		Angkatan: angkatan,
		Search:   c.Query("search"),
		Sort:     c.Query("sort", "nama"),
	}

	students, total, err := s.repo.List(c.Context(), query)
	if err != nil {
		return err
	}

	lastPage := 0
	if total > 0 {
		lastPage = (total + perPage - 1) / perPage
	}

	meta := model.Meta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    lastPage,
	}

	return helper.SuccessList(c, "Data mahasiswa berhasil diambil", students, meta)
}

// Endpoint 4: POST /api/v1/students (Admin only)
func (s *StudentService) CreateStudent(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}

	if valErrors := helper.ValidateStruct(req); len(valErrors) > 0 {
		return helper.Validation(valErrors)
	}

	// Validate angkatan <= current year
	currentYear := time.Now().Year()
	if req.Angkatan > currentYear {
		return helper.ValidationField("angkatan", "Angkatan tidak boleh melebihi tahun berjalan")
	}

	// Generate default password hash = hash(NIM)
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.NIM), bcrypt.DefaultCost)
	if err != nil {
		return helper.Internal(err)
	}

	student, err := s.repo.CreateWithUserInTx(c.Context(), req, string(passHash))
	if err != nil {
		return err
	}

	c.Set("Location", "/api/v1/students/"+strconv.Itoa(student.ID))
	return helper.Created(c, "Data mahasiswa berhasil ditambahkan", student)
}

// Endpoint 5: GET /api/v1/students/:id (Admin, Mahasiswa data sendiri)
func (s *StudentService) GetStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return err
	}

	// Authorization check: if role is mahasiswa, ensure it's their own data
	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	if authUser.Role == model.RoleMahasiswa {
		if student.UserID != authUser.UserID {
			return helper.Forbidden("tidak memiliki akses ke data mahasiswa lain")
		}
	}

	// Fetch enrolled courses and calculate SKS
	enrolledCourses, totalSKS, err := s.repo.GetEnrolledCourses(c.Context(), student.ID)
	if err != nil {
		return err
	}

	batasSKS := model.CalculateBatasSKS(student.IPKTerakhir)

	resp := model.StudentDetailResponse{
		ID:          student.ID,
		NIM:         student.NIM,
		Nama:        student.Nama,
		Prodi:       student.Prodi,
		Angkatan:    student.Angkatan,
		IPKTerakhir: student.IPKTerakhir,
		TotalSKS:    totalSKS,
		BatasSKS:    batasSKS,
		MataKuliah:  enrolledCourses,
	}

	return helper.Success(c, fiber.StatusOK, "Detail mahasiswa berhasil diambil", resp)
}

// Endpoint 6: PUT /api/v1/students/:id (Admin only)
func (s *StudentService) UpdateStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}

	if valErrors := helper.ValidateStruct(req); len(valErrors) > 0 {
		return helper.Validation(valErrors)
	}

	currentYear := time.Now().Year()
	if req.Angkatan > currentYear {
		return helper.ValidationField("angkatan", "Angkatan tidak boleh melebihi tahun berjalan")
	}

	student, err := s.repo.Update(c.Context(), id, req)
	if err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusOK, "Data mahasiswa berhasil diperbarui", student)
}

// Endpoint 7: DELETE /api/v1/students/:id (Admin only - soft delete)
func (s *StudentService) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	if err := s.repo.SoftDelete(c.Context(), id); err != nil {
		return err
	}

	return helper.NoContent(c)
}
