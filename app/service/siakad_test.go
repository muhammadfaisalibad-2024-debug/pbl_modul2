package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestApp(t *testing.T) (*fiber.App, *pgxpool.Pool, *helper.JWTManager, string, string, int, int) {
	cfg := config.Load()
	db, err := database.NewDBPool(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	jwtManager := helper.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		15*time.Minute,
	)

	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, studentRepo)
	authService := service.NewAuthService(
		userRepo,
		studentRepo,
		tokenRepo,
		jwtManager,
		7*24*time.Hour,
	)

	logger := config.NewLogger()
	app := config.NewApp(
		db,
		studentService,
		courseService,
		enrollmentService,
		authService,
		jwtManager,
		logger,
	)

	// Fetch admin and student1 info
	adminUser, err := userRepo.FindByEmail(context.Background(), "admin@siakad.ac.id")
	if err != nil {
		t.Fatalf("Admin user not found: %v", err)
	}
	adminToken, _ := jwtManager.GenerateAccess(*adminUser)

	mhsUser, err := userRepo.FindByEmail(context.Background(), "rina.putri@siakad.ac.id")
	if err != nil {
		t.Fatalf("Student user not found: %v", err)
	}
	mhsToken, _ := jwtManager.GenerateAccess(*mhsUser)

	mhsStudent, err := studentRepo.FindByUserID(context.Background(), mhsUser.ID)
	if err != nil {
		t.Fatalf("Student profile not found: %v", err)
	}

	return app, db, jwtManager, adminToken, mhsToken, adminUser.ID, mhsStudent.ID
}

// Test Business Rule 1: Batas SKS berdasarkan IPK
func TestBatasSKSCalculation(t *testing.T) {
	tests := []struct {
		ipk      float64
		expected int
	}{
		{4.00, 24},
		{3.50, 24},
		{3.00, 24},
		{2.99, 21},
		{2.75, 21},
		{2.50, 21},
		{2.49, 18},
		{2.00, 18},
		{0.00, 18},
	}

	for _, tc := range tests {
		got := model.CalculateBatasSKS(tc.ipk)
		if got != tc.expected {
			t.Errorf("For IPK %.2f expected %d SKS, got %d", tc.ipk, tc.expected, got)
		}
	}
}

// Test Endpoint 1: POST /api/v1/auth/login
func TestEndpoint1_Login(t *testing.T) {
	app, db, _, _, _, _, _ := setupTestApp(t)
	defer db.Close()

	// 1. Success Admin Login (200)
	body, _ := json.Marshal(map[string]string{
		"email":    "admin@siakad.ac.id",
		"password": "Admin123!",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Admin login failed, status: %d, err: %v", resp.StatusCode, err)
	}

	// 2. Success Student Login (200)
	body, _ = json.Marshal(map[string]string{
		"email":    "rina.putri@siakad.ac.id",
		"password": "187221000001",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Student login failed, status: %d, err: %v", resp.StatusCode, err)
	}

	// 3. Wrong Password (401)
	body, _ = json.Marshal(map[string]string{
		"email":    "admin@siakad.ac.id",
		"password": "wrongpassword123",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 401 {
		t.Errorf("Expected status 401 for wrong credentials, got %d", resp.StatusCode)
	}

	// 4. Invalid Format (422)
	body, _ = json.Marshal(map[string]string{
		"email":    "invalid-email-format",
		"password": "short",
	})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 422 {
		t.Errorf("Expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

// Test Endpoint 2: GET /api/v1/auth/me
func TestEndpoint2_AuthMe(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, _ := setupTestApp(t)
	defer db.Close()

	// 1. Admin Me (200)
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Admin Me failed, status: %d", resp.StatusCode)
	}

	// 2. Student Me (200, includes student profile)
	req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, err = app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Student Me failed, status: %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var res struct {
		Data model.MeResponse `json:"data"`
	}
	_ = json.Unmarshal(respBody, &res)
	if res.Data.Student == nil || res.Data.Student.NIM != "187221000001" {
		t.Errorf("Expected student data attached, got: %s", string(respBody))
	}

	// 3. Unauthorized without token (401)
	req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 401 {
		t.Errorf("Expected 401 for unauthorized /me, got %d", resp.StatusCode)
	}
}

// Test Endpoint 3: GET /api/v1/students (Admin only)
func TestEndpoint3_GetStudents(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, _ := setupTestApp(t)
	defer db.Close()

	// 1. Admin access with pagination & search (200)
	req := httptest.NewRequest("GET", "/api/v1/students?page=1&per_page=10&search=Rina", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("GET /students failed, status: %d", resp.StatusCode)
	}

	// 2. Mahasiswa access (403)
	req = httptest.NewRequest("GET", "/api/v1/students", nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for mahasiswa accessing GET /students, got %d", resp.StatusCode)
	}
}

// Test Endpoint 4: POST /api/v1/students (Admin only)
func TestEndpoint4_CreateStudent(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, _ := setupTestApp(t)
	defer db.Close()

	testNIM := "187221000999"
	testEmail := "test.student999@siakad.ac.id"

	// Cleanup test student if exists
	_, _ = db.Exec(context.Background(), "DELETE FROM students WHERE nim = $1", testNIM)
	_, _ = db.Exec(context.Background(), "DELETE FROM users WHERE email = $1", testEmail)

	// 1. Mahasiswa forbidden (403)
	body, _ := json.Marshal(map[string]interface{}{
		"nim":      testNIM,
		"nama":     "Test Student",
		"email":    testEmail,
		"prodi":    "Sistem Informasi",
		"angkatan": 2024,
	})
	req := httptest.NewRequest("POST", "/api/v1/students", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ := app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for mahasiswa creating student, got %d", resp.StatusCode)
	}

	// 2. Admin creates student (201)
	req = httptest.NewRequest("POST", "/api/v1/students", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 201 {
		t.Fatalf("Admin create student failed, status: %d", resp.StatusCode)
	}

	// 3. Duplicate NIM returns 422
	req = httptest.NewRequest("POST", "/api/v1/students", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 422 {
		t.Errorf("Expected 422 for duplicate student, got %d", resp.StatusCode)
	}

	// Cleanup
	_, _ = db.Exec(context.Background(), "DELETE FROM students WHERE nim = $1", testNIM)
	_, _ = db.Exec(context.Background(), "DELETE FROM users WHERE email = $1", testEmail)
}

// Test Endpoint 5: GET /api/v1/students/:id (Admin & Mahasiswa own data)
func TestEndpoint5_GetStudentDetail(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, mhsStudentID := setupTestApp(t)
	defer db.Close()

	// 1. Student accesses own detail (200)
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/students/%d", mhsStudentID), nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Student access own detail failed, status: %d", resp.StatusCode)
	}

	// 2. Student accesses other student detail (403)
	// Find another student ID
	var otherID int
	_ = db.QueryRow(context.Background(), "SELECT id FROM students WHERE id != $1 AND deleted_at IS NULL LIMIT 1", mhsStudentID).Scan(&otherID)

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/students/%d", otherID), nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for student accessing other student, got %d", resp.StatusCode)
	}

	// 3. Admin accesses any student detail (200)
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/students/%d", otherID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Admin access student detail failed, status: %d", resp.StatusCode)
	}

	// 4. Not found (404)
	req = httptest.NewRequest("GET", "/api/v1/students/999999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404 for nonexistent student, got %d", resp.StatusCode)
	}
}

// Test Endpoint 6: PUT /api/v1/students/:id (Admin only)
func TestEndpoint6_UpdateStudent(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, mhsStudentID := setupTestApp(t)
	defer db.Close()

	body, _ := json.Marshal(map[string]interface{}{
		"nama":         "Rina Putri Updated",
		"prodi":        "Sistem Informasi",
		"angkatan":     2022,
		"ipk_terakhir": 3.75,
	})

	// 1. Mahasiswa forbidden (403)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/students/%d", mhsStudentID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ := app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for student updating student, got %d", resp.StatusCode)
	}

	// 2. Admin update (200)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/students/%d", mhsStudentID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Admin update student failed, status: %d", resp.StatusCode)
	}
}

// Test Endpoint 7: DELETE /api/v1/students/:id (Soft Delete)
func TestEndpoint7_DeleteStudent(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, _ := setupTestApp(t)
	defer db.Close()

	// Create a temp student for deletion test
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)

	reqData := model.CreateStudentRequest{
		NIM:      "187221000888",
		Nama:     "Temp Student",
		Email:    "temp.student888@siakad.ac.id",
		Prodi:    "Informatika",
		Angkatan: 2024,
	}
	createdStudent, err := studentRepo.CreateWithUserInTx(context.Background(), reqData, "dummyhash")
	if err != nil {
		t.Fatalf("Failed to create temp student: %v", err)
	}

	// 1. Mahasiswa forbidden (403)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/students/%d", createdStudent.ID), nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ := app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for student deleting student, got %d", resp.StatusCode)
	}

	// 2. Admin delete (204)
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/students/%d", createdStudent.ID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = app.Test(req, 5000)
	if err != nil || resp.StatusCode != 204 {
		t.Fatalf("Admin delete student failed, status: %d", resp.StatusCode)
	}

	// 3. Cannot get deleted student (404)
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/students/%d", createdStudent.ID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404 for deleted student, got %d", resp.StatusCode)
	}

	// Cleanup
	_ = userRepo
	_, _ = db.Exec(context.Background(), "DELETE FROM students WHERE id = $1", createdStudent.ID)
	_, _ = db.Exec(context.Background(), "DELETE FROM users WHERE email = $1", reqData.Email)
}

// Test Endpoint 8: GET /api/v1/courses (All roles)
func TestEndpoint8_GetCourses(t *testing.T) {
	app, db, _, _, mhsToken, _, _ := setupTestApp(t)
	defer db.Close()

	// 1. Get all courses with query parameters (200)
	req := httptest.NewRequest("GET", "/api/v1/courses?semester=3&available=true", nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("GET /courses failed, status: %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var res struct {
		Data []model.CourseListItem `json:"data"`
	}
	_ = json.Unmarshal(respBody, &res)
	if len(res.Data) == 0 {
		t.Errorf("Expected courses returned, got empty")
	}
}

// Test Endpoint 9 & 10: POST /api/v1/enrollments & DELETE /api/v1/enrollments/:id
func TestEndpoint9And10_Enrollments(t *testing.T) {
	app, db, _, adminToken, mhsToken, _, mhsStudentID := setupTestApp(t)
	defer db.Close()

	// Get a course ID
	var courseID int
	_ = db.QueryRow(context.Background(), "SELECT id FROM courses WHERE kode_mk = 'MK001'").Scan(&courseID)

	// Clean any previous test enrollment
	_, _ = db.Exec(context.Background(), "DELETE FROM enrollments WHERE student_id = $1", mhsStudentID)

	// 1. Admin forbidden (403)
	body, _ := json.Marshal(map[string]interface{}{
		"course_id":      courseID,
		"tahun_akademik": "2026/2027-Ganjil",
	})
	req := httptest.NewRequest("POST", "/api/v1/enrollments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, _ := app.Test(req, 5000)
	if resp.StatusCode != 403 {
		t.Errorf("Expected 403 for admin taking course, got %d", resp.StatusCode)
	}

	// 2. Mahasiswa takes course (201)
	req = httptest.NewRequest("POST", "/api/v1/enrollments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, err := app.Test(req, 5000)
	if err != nil || resp.StatusCode != 201 {
		t.Fatalf("Student take course failed, status: %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var enrRes struct {
		Data model.Enrollment `json:"data"`
	}
	_ = json.Unmarshal(respBody, &enrRes)
	enrollmentID := enrRes.Data.ID

	// 3. Duplicate enrollment in same academic year (409 Conflict)
	req = httptest.NewRequest("POST", "/api/v1/enrollments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 409 {
		t.Errorf("Expected 409 for duplicate enrollment, got %d", resp.StatusCode)
	}

	// 4. Mahasiswa deletes own enrollment (204)
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/enrollments/%d", enrollmentID), nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, err = app.Test(req, 5000)
	if err != nil || resp.StatusCode != 204 {
		t.Fatalf("Student cancel enrollment failed, status: %d", resp.StatusCode)
	}

	// 5. Delete nonexistent enrollment (404)
	req = httptest.NewRequest("DELETE", "/api/v1/enrollments/999999", nil)
	req.Header.Set("Authorization", "Bearer "+mhsToken)
	resp, _ = app.Test(req, 5000)
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404 for deleting nonexistent enrollment, got %d", resp.StatusCode)
	}
}
