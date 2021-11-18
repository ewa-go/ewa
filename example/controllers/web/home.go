package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Home struct {
	utils.Navbar
}

func (h *Home) Get(route *ewa.Route) {
	route.Handler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		//h.Navbar = utils.NewNavbar("", identity.User)
		return ctx.Render("home", h, "layouts/base")
	}
	/*route.Handler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		//h.Navbar = utils.NewNavbar("", identity.User)
		return ctx.Render("home", h, "layouts/base")
	}*/
}

func (h *Home) Post(route *ewa.Route) {
	route.Empty()
}
