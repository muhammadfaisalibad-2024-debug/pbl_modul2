package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type PortfolioService struct {
	portfolios repository.PortfolioRepository
	students   *repository.StudentRepository
	perms      *helper.PermissionSet
}

func NewPortfolioService(
	portfolios repository.PortfolioRepository,
	students *repository.StudentRepository,
	perms *helper.PermissionSet,
) *PortfolioService {
	return &PortfolioService{
		portfolios: portfolios,
		students:   students,
		perms:      perms,
	}
}

// GET /api/v1/students/:id/portfolio
func (s *PortfolioService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	// 1. Pastikan mahasiswa terdaftar di PostgreSQL
	if _, err := s.students.FindByID(ctx, studentID); err != nil {
		return helper.NotFound("mahasiswa tidak ditemukan di database")
	}

	// 2. Ambil portofolio dari MongoDB
	portfolio, err := s.portfolios.FindByStudent(ctx, studentID)
	if err != nil {
		return translatePortfolioError(err)
	}

	return helper.Success(c, fiber.StatusOK, "portofolio berhasil diambil", portfolio)
}

// PUT /api/v1/students/:id/portfolio
func (s *PortfolioService) Put(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	var req model.PutPortfolioRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)

	for i := range req.Projects {
		req.Projects[i].Title = strings.TrimSpace(req.Projects[i].Title)
		req.Projects[i].URL = strings.TrimSpace(req.Projects[i].URL)
		req.Projects[i].Technologies = NormalizeTags(req.Projects[i].Technologies)
	}

	for i := range req.Certificates {
		req.Certificates[i].Name = strings.TrimSpace(req.Certificates[i].Name)
		req.Certificates[i].Issuer = strings.TrimSpace(req.Certificates[i].Issuer)
		req.Certificates[i].Date = strings.TrimSpace(req.Certificates[i].Date)
	}

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	// 1. Pastikan mahasiswa ada di PostgreSQL
	student, err := s.students.FindByID(ctx, studentID)
	if err != nil {
		return helper.NotFound("mahasiswa tidak ditemukan di database")
	}

	// 2. Otorisasi: Mahasiswa pemilik data atau Admin dengan permission portfolio:manage
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	canManage := s.perms.Can(current.Role, "portfolio:manage")
	isOwner := (student.UserID == current.UserID)
	if !canManage && !isOwner {
		return helper.Forbidden("hanya pemilik portofolio atau admin yang dapat mengubah data")
	}

	// 3. Konversi struct model
	projects := make([]model.Project, 0, len(req.Projects))
	for _, p := range req.Projects {
		projects = append(projects, model.Project{
			Title:        p.Title,
			Technologies: p.Technologies,
			URL:          p.URL,
			Year:         p.Year,
		})
	}

	certs := make([]model.Certificate, 0, len(req.Certificates))
	for _, certInput := range req.Certificates {
		parsedDate, parseErr := parseCertificateDate(certInput.Date)
		if parseErr != nil {
			return helper.ValidationField("date", "format tanggal sertifikat tidak valid, gunakan YYYY-MM-DD")
		}
		certs = append(certs, model.Certificate{
			Name:   certInput.Name,
			Issuer: certInput.Issuer,
			Date:   parsedDate,
		})
	}

	saved, err := s.portfolios.Replace(ctx, model.Portfolio{
		StudentID:    studentID,
		Headline:     req.Headline,
		Bio:          req.Bio,
		Projects:     projects,
		Certificates: certs,
	})
	if err != nil {
		return translatePortfolioError(err)
	}

	// Berdasarkan keputusan user, PUT selalu mengembalikan 200 OK
	return helper.Success(c, fiber.StatusOK, "portofolio berhasil disimpan", saved)
}

// POST /api/v1/students/:id/portfolio/projects
func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	var req model.AddProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.URL = strings.TrimSpace(req.URL)
	req.Technologies = NormalizeTags(req.Technologies)

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	// 1. Cek keberadaan mahasiswa di PostgreSQL
	student, err := s.students.FindByID(ctx, studentID)
	if err != nil {
		return helper.NotFound("mahasiswa tidak ditemukan di database")
	}

	// 2. Otorisasi
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	canManage := s.perms.Can(current.Role, "portfolio:manage")
	isOwner := (student.UserID == current.UserID)
	if !canManage && !isOwner {
		return helper.Forbidden("hanya pemilik portofolio atau admin yang dapat mengubah data")
	}

	// 3. Tambah proyek secara atomik
	updated, err := s.portfolios.AddProject(ctx, studentID, model.Project{
		Title:        req.Title,
		Technologies: req.Technologies,
		URL:          req.URL,
		Year:         req.Year,
	})
	if err != nil {
		return translatePortfolioError(err)
	}

	location := "/api/v1/students/" + strconv.Itoa(studentID) + "/portfolio"
	return helper.Created(c, "proyek berhasil ditambahkan ke portofolio", updated, location)
}

// DELETE /api/v1/students/:id/portfolio/projects/:projectId
func (s *PortfolioService) RemoveProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	projectID := c.Params("projectId")
	if projectID == "" {
		return helper.BadRequest("projectId harus diisi")
	}

	// 1. Cek keberadaan mahasiswa di PostgreSQL
	student, err := s.students.FindByID(ctx, studentID)
	if err != nil {
		return helper.NotFound("mahasiswa tidak ditemukan di database")
	}

	// 2. Otorisasi
	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	canManage := s.perms.Can(current.Role, "portfolio:manage")
	isOwner := (student.UserID == current.UserID)
	if !canManage && !isOwner {
		return helper.Forbidden("hanya pemilik portofolio atau admin yang dapat mengubah data")
	}

	// 3. Hapus proyek dari MongoDB
	if err := s.portfolios.RemoveProject(ctx, studentID, projectID); err != nil {
		return translatePortfolioError(err)
	}

	return helper.NoContent(c)
}

// GET /api/v1/portfolios/stats/technologies
func (s *PortfolioService) TopTechnologies(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	limit := c.QueryInt("limit", 10)
	stats, err := s.portfolios.TopTechnologies(ctx, limit)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "statistik 10 teknologi terpopuler berhasil diambil", stats)
}

func parseCertificateDate(dateStr string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t.UTC().Truncate(time.Millisecond), nil
	}
	return time.Time{}, errors.New("format tanggal tidak valid")
}

func translatePortfolioError(err error) error {
	switch {
	case errors.Is(err, repository.ErrInvalidID):
		return helper.BadRequest("format id tidak valid")
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("portofolio tidak ditemukan")
	case errors.Is(err, repository.ErrProjectNotFound):
		return helper.NotFound("proyek tidak ditemukan dalam portofolio")
	case errors.Is(err, repository.ErrProjectLimit):
		return helper.NewError(fiber.StatusConflict, "PROJECT_LIMIT", "jumlah proyek sudah mencapai batas maksimum (20 proyek)")
	default:
		return helper.Internal(err)
	}
}
