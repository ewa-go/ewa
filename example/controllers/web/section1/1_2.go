package section1

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Section_1_2 struct {
	utils.Navbar
	Title string
}

func (h *Section_1_2) Get(route *ewa.Route) {
	route.SetDescription("Страница Home.html").Session(true)
	route.WebHandler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		h.Navbar = utils.NewNavbar("section_1_2", identity.User)
		h.Title = "Привет раздел 1.2"
		return ctx.Render("section1/1_2", h, "layouts/base")
	}
}

func (h *Section_1_2) Post(route *ewa.Route) {
	route.Handler = nil
}
