package helper

import (
	"api-students/app/model"
	"encoding/csv"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

const FormatCSV = "text/csv"

func Negotiate(c *fiber.Ctx, offered ...string) (string, error) {
	accept := strings.TrimSpace(c.Get(fiber.HeaderAccept))
	if accept == "" || accept == "*/*" {
		return offered[0], nil
	}
	chosen := c.Accepts(offered...)
	if chosen == "" {
		return "", NotAcceptable("format yang diminta tidak tersedia")
	}
	return chosen, nil
}
func WriteUsersCSV(c *fiber.Ctx, users []model.User) error {
	c.Set(fiber.HeaderContentType, FormatCSV+"; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="users.csv"`)
	var b strings.Builder
	w := csv.NewWriter(&b)
	_ = w.Write([]string{"id", "username", "email", "role", "is_active", "created_at"})
	for _, u := range users {
		if err := w.Write([]string{strconv.Itoa(u.ID), u.Username, u.Email, u.Role, strconv.FormatBool(u.IsActive), u.CreatedAt.UTC().Format(time.RFC3339)}); err != nil {
			return Internal(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return Internal(err)
	}
	return c.SendString(b.String())
}

func WriteStudentsCSV(c *fiber.Ctx, students []model.Student) error {
	c.Set(fiber.HeaderContentType, FormatCSV+"; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="students.csv"`)
	var b strings.Builder
	w := csv.NewWriter(&b)
	_ = w.Write([]string{"id", "nim", "name", "grade", "is_active", "owner_id", "created_at"})
	for _, s := range students {
		if err := w.Write([]string{strconv.Itoa(s.ID), s.NIM, s.Name, strconv.FormatFloat(s.Grade, 'f', -1, 64), strconv.FormatBool(s.IsActive), strconv.Itoa(s.OwnerID), s.CreatedAt.UTC().Format(time.RFC3339)}); err != nil {
			return Internal(err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return Internal(err)
	}
	return c.SendString(b.String())
}
