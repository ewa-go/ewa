package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Home struct {
	utils.NavBar
}

func (h *Home) Get(route *ewa.Route) {
	route.SetDescription("Страница Home.html").Session()
	route.Handler = func(ctx *fiber.Ctx) error {
		h.NavBar = utils.GetNavBar("")
		return ctx.Render("home", h, "layouts/base")
	}
}

func (h *Home) Post(route *ewa.Route) {
	route.Handler = nil
}
