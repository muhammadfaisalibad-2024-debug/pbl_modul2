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

type UserService struct {
	repo  repository.UserRepository
	perms *helper.PermissionSet
}

func NewUserService(repo repository.UserRepository, perms *helper.PermissionSet) *UserService {
	return &UserService{repo: repo, perms: perms}
}

func (s *UserService) List(c *fiber.Ctx) error {
	format, err := helper.Negotiate(c, fiber.MIMEApplicationJSON, helper.FormatCSV)
	if err != nil {
		return err
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	query := model.CursorQuery{Limit: limit, Search: strings.TrimSpace(c.Query("search"))}
	if raw := strings.TrimSpace(c.Query("cursor")); raw != "" {
		cursor, e := helper.DecodeCursor(raw)
		if e != nil {
			return helper.BadRequest("cursor tidak valid")
		}
		query.After = &cursor
	}
	users, err := s.repo.FindAfterCursor(c.Context(), query)
	if err != nil {
		return helper.Internal(err)
	}
	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}
	meta := model.CursorMeta{Limit: limit, HasMore: hasMore}
	if hasMore {
		last := users[len(users)-1]
		meta.NextCursor = helper.EncodeCursor(last.CreatedAt, last.ID)
	}
	if format == helper.FormatCSV {
		return helper.WriteUsersCSV(c, users)
	}
	return helper.SuccessCursor(c, "daftar user berhasil diambil", users, meta)
}

func (s *UserService) Create(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	req.Username, req.Email = strings.TrimSpace(req.Username), strings.TrimSpace(req.Email)
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}
	password, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal memproses password")
	}
	user, err := s.repo.Create(c.Context(), model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: password,
		Role:     "user",
		IsActive: true,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, "username atau email sudah dipakai")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal membuat user")
	}
	c.Location("/api/v1/users/" + strconv.Itoa(user.ID))
	return helper.Created(c, "user berhasil dibuat", user)
}

func (s *UserService) Get(c *fiber.Ctx) error {
	current, id, ok := currentAndID(c)
	if !ok {
		return nil
	}
	if !CanAccessUser(current, id, s.perms, "user:read:any") {
		return helper.Fail(c, fiber.StatusForbidden, "tidak berhak mengakses data user lain")
	}
	user, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return userError(c, err, "gagal mengambil data user")
	}
	return helper.Success(c, fiber.StatusOK, "user ditemukan", user)
}

func (s *UserService) Replace(c *fiber.Ctx) error {
	current, id, ok := currentAndID(c)
	if !ok {
		return nil
	}
	if !CanAccessUser(current, id, s.perms, "user:update:any") {
		return helper.Fail(c, fiber.StatusForbidden, "tidak berhak mengubah data user lain")
	}
	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	req.Username, req.Email = strings.TrimSpace(req.Username), strings.TrimSpace(req.Email)
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}
	user, err := s.repo.Update(c.Context(), id, req)
	if err != nil {
		return userError(c, err, "gagal mengubah user")
	}
	return helper.Success(c, fiber.StatusOK, "user berhasil diubah", user)
}

func (s *UserService) Patch(c *fiber.Ctx) error {
	current, id, ok := currentAndID(c)
	if !ok {
		return nil
	}
	if !CanAccessUser(current, id, s.perms, "user:update:any") {
		return helper.Fail(c, fiber.StatusForbidden, "tidak berhak mengubah data user lain")
	}
	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	if req.Username == nil && req.Email == nil && req.IsActive == nil {
		return helper.BadRequest("tidak ada field yang diubah")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}
	user, err := s.repo.Patch(c.Context(), id, req)
	if err != nil {
		return userError(c, err, "gagal mengubah user")
	}
	return helper.Success(c, fiber.StatusOK, "user berhasil diubah", user)
}

func (s *UserService) AssignRole(c *fiber.Ctx) error {
	current, id, ok := currentAndID(c)
	if !ok {
		return nil
	}
	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	if errs := ValidateAssignRole(current, id, req, s.perms); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"success": false, "message": "validation failed", "errors": errs})
	}
	user, err := s.repo.UpdateRole(c.Context(), id, strings.TrimSpace(req.Role))
	if err != nil {
		return userError(c, err, "gagal mengubah role user")
	}
	return helper.Success(c, fiber.StatusOK, "role user berhasil diubah", user)
}

func (s *UserService) Delete(c *fiber.Ctx) error {
	current, id, ok := currentAndID(c)
	if !ok {
		return nil
	}
	if current.UserID == id {
		return helper.Fail(c, fiber.StatusForbidden, "tidak boleh menghapus akun sendiri")
	}
	if err := s.repo.Delete(c.Context(), id); err != nil {
		return userError(c, err, "gagal menghapus user")
	}
	return helper.NoContent(c)
}

func currentAndID(c *fiber.Ctx) (model.AuthUser, int, bool) {
	current, ok := helper.CurrentUser(c)
	if !ok {
		_ = helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
		return model.AuthUser{}, 0, false
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		_ = helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
		return model.AuthUser{}, 0, false
	}
	return current, id, true
}

func userError(c *fiber.Ctx, err error, fallback string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return helper.Fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	return helper.Fail(c, fiber.StatusInternalServerError, fallback)
}
