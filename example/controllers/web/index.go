package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Index struct{}

func (a *Index) Get(route *ewa.Route) {
	route.SetDescription("Страница Index.html")
	route.Handler = func(ctx *fiber.Ctx) error {
		return ctx.Render("index", nil)
	}
}

func (a *Index) Post(route *ewa.Route) {

}
