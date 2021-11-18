package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/egovorukhin/egowebapi/swagger"
	"github.com/gofiber/fiber/v2"
)

type Swagger struct {
	utils.Navbar
}

func (s *Swagger) Get(route *ewa.Route) {
	route.SwaggerHandler = func(ctx *fiber.Ctx, swagger *swagger.Swagger) error {
		return ctx.JSON(swagger)
	}
}

func (s *Swagger) Post(route *ewa.Route) {
	route.Empty()
}