package service

import (
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	studentRepo *repository.StudentRepository
	tokenRepo   *repository.TokenRepository
	jwtManager  *helper.JWTManager
	refreshTTL  time.Duration
}

func NewAuthService(
	userRepo *repository.UserRepository,
	studentRepo *repository.StudentRepository,
	tokenRepo *repository.TokenRepository,
	jwtManager *helper.JWTManager,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		studentRepo: studentRepo,
		tokenRepo:   tokenRepo,
		jwtManager:  jwtManager,
		refreshTTL:  refreshTTL,
	}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}

	if valErrors := helper.ValidateStruct(req); len(valErrors) > 0 {
		return helper.Validation(valErrors)
	}

	// Find user
	user, err := s.userRepo.FindByEmail(c.Context(), req.Email)
	if err != nil {
		return helper.Unauthorized("email atau password salah")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return helper.Unauthorized("email atau password salah")
	}

	// Check if student is soft-deleted
	if user.Role == model.RoleMahasiswa {
		student, err := s.studentRepo.FindByUserID(c.Context(), user.ID)
		if err != nil || (student != nil && student.DeletedAt != nil) {
			return helper.Unauthorized("akun mahasiswa telah dinonaktifkan atau dihapus")
		}
	}

	// Generate access token
	accessToken, err := s.jwtManager.GenerateAccess(*user)
	if err != nil {
		return helper.Internal(err)
	}

	// Generate refresh token
	refreshToken, err := helper.GenerateSecureToken()
	if err == nil && s.tokenRepo != nil {
		tokenHash := helper.HashToken(refreshToken)
		_ = s.tokenRepo.Save(c.Context(), user.ID, tokenHash, time.Now().Add(s.refreshTTL))
	}

	return helper.Success(c, fiber.StatusOK, "Login berhasil", fiber.Map{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(s.jwtManager.AccessTTL().Seconds()),
		"user": model.UserData{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	user, err := s.userRepo.FindByID(c.Context(), authUser.UserID)
	if err != nil {
		return helper.Unauthorized("user tidak ditemukan")
	}

	resp := model.MeResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	if user.Role == model.RoleMahasiswa {
		student, err := s.studentRepo.FindByUserID(c.Context(), user.ID)
		if err == nil && student != nil {
			resp.Student = &model.StudentProfile{
				NIM:      student.NIM,
				Nama:     student.Nama,
				Prodi:    student.Prodi,
				Angkatan: student.Angkatan,
			}
		}
	}

	return helper.Success(c, fiber.StatusOK, "Profil pengguna berhasil diambil", resp)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}
	if req.RefreshToken == "" {
		return helper.BadRequest("refresh_token wajib diisi")
	}

	tokenHash := helper.HashToken(req.RefreshToken)
	storedToken, err := s.tokenRepo.FindValid(c.Context(), tokenHash)
	if err != nil {
		return helper.Unauthorized("refresh token tidak valid atau sudah kedaluwarsa")
	}

	user, err := s.userRepo.FindByID(c.Context(), storedToken.UserID)
	if err != nil {
		return helper.Unauthorized("user tidak ditemukan")
	}

	_ = s.tokenRepo.Revoke(c.Context(), storedToken.ID)

	newAccessToken, err := s.jwtManager.GenerateAccess(*user)
	if err != nil {
		return helper.Internal(err)
	}

	newRefreshToken, err := helper.GenerateSecureToken()
	if err != nil {
		return helper.Internal(err)
	}

	newTokenHash := helper.HashToken(newRefreshToken)
	if err := s.tokenRepo.Save(c.Context(), user.ID, newTokenHash, time.Now().Add(s.refreshTTL)); err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Token berhasil diperbarui", model.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.AccessTTL().Seconds()),
	})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		tokenHash := helper.HashToken(req.RefreshToken)
		_ = s.tokenRepo.RevokeByHash(c.Context(), tokenHash)
	}
	return helper.Success(c, fiber.StatusOK, "Logout berhasil", nil)
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Format JSON tidak valid")
	}
	if valErrors := helper.ValidateStruct(req); len(valErrors) > 0 {
		return helper.Validation(valErrors)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.Internal(err)
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Role:     model.RoleMahasiswa,
	}

	if err := s.userRepo.Create(c.Context(), &user); err != nil {
		return err
	}

	return helper.Created(c, "Registrasi berhasil", model.UserData{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}
