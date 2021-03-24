package web

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber"
)

type Index struct {
	*egowebapi.Controller
}

func NewIndex(path string) Index {

	a := Index{
		Controller: egowebapi.NewController("Index", "Страница Index.html"),
	}

	path = a.CheckPath(path, a)

	a.SetRoutes(
		egowebapi.NewRoute("GET", path, a.Get),
		egowebapi.NewRoute("POST", path, a.Post),
	)

	return a
}

func (a Index) Get(c *fiber.Ctx) {
	_ = c.Render("index", nil)
}

func (a Index) Post(c *fiber.Ctx) {

}
