package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Index struct{}

func (a *Index) Get() *ewa.Route {
	return &ewa.Route{
		Params: nil,
		Handler: func(ctx *fiber.Ctx) error {
			return ctx.Render("index", nil)
		},
	}
}

func (a *Index) Post() *ewa.Route {
	return nil
}
