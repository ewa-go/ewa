package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Swagger struct {
	utils.Navbar
}

func (s *Swagger) Get(route *ewa.Route) {
	route.Handler = func(ctx *fiber.Ctx, swagger *ewa.Swagger) error {
		return ctx.JSON(swagger)
	}
}
