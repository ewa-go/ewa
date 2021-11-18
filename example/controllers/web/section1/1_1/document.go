package __1

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Document struct {
	utils.Navbar
	Title string
}

func (h *Document) Get(route *ewa.Route) {
	route.SetDescription("Страница Home.html").Session().Permission()
	route.Handler = func(ctx *fiber.Ctx, identity *ewa.Identity) error {
		h.Navbar = utils.NewNavbar("section1/1_1/document", identity.User)
		h.Title = "Привет раздел 1.1"
		return ctx.Render("section1/1_1/document", h, "layouts/base")
	}
}

func (h *Document) Post(route *ewa.Route) {
	route.Handler = nil
}
