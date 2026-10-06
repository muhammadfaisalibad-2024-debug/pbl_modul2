package helper

import (
	"strings"

	"api-students/app/model"
	"github.com/gofiber/fiber/v2"
)

func ParseListQuery(c *fiber.Ctx) model.StudentListQuery {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}
	if perPage > 50 {
		perPage = 50
	}

	search := strings.TrimSpace(c.Query("search"))
	sort := c.Query("sort", "nama")
	prodi := c.Query("prodi")
	angkatan := c.QueryInt("angkatan", 0)

	return model.StudentListQuery{
		Page:     page,
		PerPage:  perPage,
		Search:   search,
		Sort:     sort,
		Prodi:    prodi,
		Angkatan: angkatan,
	}
}
