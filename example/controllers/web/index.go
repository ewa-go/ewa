package web

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Index struct {
}

func (a *Index) Get() *egowebapi.Route {
	return &egowebapi.Route{
		Path: egowebapi.AddPath(""),
		Handler: func(ctx *fiber.Ctx) error {
			return ctx.Render("/web/index", nil)
		},
	}
}

func (a *Index) Post() *egowebapi.Route {
	return nil
}
