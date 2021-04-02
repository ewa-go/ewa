package web

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber"
)

type Index struct {
}

func (a *Index) Get() *egowebapi.Route {
	return &egowebapi.Route{
		Path: egowebapi.AddPath(""),
		Handler: func(ctx *fiber.Ctx) {
			_ = ctx.Render("index", nil)
		},
	}
}

func (a *Index) Post() *egowebapi.Route {
	return nil
}
