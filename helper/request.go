package helper

import (
	"strings"

	"api-students/app/model"
	"github.com/gofiber/fiber/v2"
)

func ParseListQuery(c *fiber.Ctx) model.ListQuery {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	search := strings.ToLower(c.Query("search"))
	sortBy := c.Query("sort", "id")
	order := strings.ToLower(c.Query("order", "asc"))
	active := c.Query("is_active")

	return model.ListQuery{
		Page:     page,
		Limit:    limit,
		Search:   search,
		SortBy:   sortBy,
		Order:    order,
		IsActive: active,
	}
}
