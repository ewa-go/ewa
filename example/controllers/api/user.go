package api

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/valyala/fasthttp"
)

type User struct {
	*egowebapi.Controller
}

func NewUser(path string) *egowebapi.Controller {

	a := User{
		Controller: &egowebapi.Controller{
			Name:        "User",
			Description: "Контроллер для пользователей",
		},
	}

	routes := egowebapi.NewRoutes(
		egowebapi.NewRoute(path, "GET", "Вернуть пользователей", a.Get),
		egowebapi.NewRoute(path, "POST", "Добавить пользователя", a.Post),
		egowebapi.NewRoute(path, "PUT", "Изменить пользователя", a.Put),
		egowebapi.NewRoute(path, "DELETE", "Удалить подльзователя", a.Delete),
	)
	a.Routes = routes

	return a.Controller
}

func (a *User) Get(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Hello, World!"))
}

func (a *User) Post(ctx *fasthttp.RequestCtx) {

}

func (a *User) Put(ctx *fasthttp.RequestCtx) {

}

func (a *User) Delete(ctx *fasthttp.RequestCtx) {

}
