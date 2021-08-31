package __1

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type List struct {
	utils.Navbar
	Title string
}

func (h *List) Get(route *ewa.Route) {
	route.SetDescription("Страница Home.html").Session(true)
	route.WebHandler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		h.Navbar = utils.NewNavbar("section1/1_1/list", identity.User)
		h.Title = "Привет раздел 1.1"
		return ctx.Render("section1/1_1/list", h, "layouts/base")
	}
}

func (h *List) Post(route *ewa.Route) {
	route.Handler = nil
}
