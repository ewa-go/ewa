package web

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/valyala/fasthttp"
)

type Index struct {
	*egowebapi.Controller
}

func NewIndex(path string) *egowebapi.Controller {

	a := Index{
		Controller: egowebapi.NewController("Index", "Страница Index.html"),
	}

	routes := egowebapi.NewRoutes(
		egowebapi.NewRoute(path, "GET", "Метод GET", a.Get),
		egowebapi.NewRoute(path, "POST", "Метод POST", a.Post),
	)
	a.Routes = routes

	return a.Controller
}

func (a *Index) Get(ctx *fasthttp.RequestCtx) {
	_ = a.View(ctx.Response.BodyWriter(), "", nil)
}

func (a *Index) Post(ctx *fasthttp.RequestCtx) {

}

