package middleware

import (
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(perms *helper.PermissionSet, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Unauthorized("belum terautentikasi")
		}
		if !perms.Can(user.Role, permission) {
			return helper.Forbidden("role " + user.Role + " tidak memiliki hak akses " + permission)
		}
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Unauthorized("belum terautentikasi")
		}
		for _, r := range roles {
			if user.Role == r {
				return c.Next()
			}
		}
		return helper.Forbidden("akses ditolak, role " + user.Role + " tidak diizinkan")
	}
}
