package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

type Route struct {
	Params       []string
	Description  string
	IsBasicAuth  bool
	IsJWT        bool
	IsSession    bool
	IsPermission bool
	Handler      Handler
	WebHandler   WebHandler
}

const sessionId = "session_id"

func (r *Route) SetHandler(handler Handler) *Route {
	r.Handler = handler
	return r
}

func (r *Route) SetParams(params ...string) *Route {
	r.Params = params
	return r
}

func (r *Route) SetDescription(s string) *Route {
	r.Description = s
	return r
}

// BasicAuth Вешаем флаг авторизации Basic
func (r *Route) BasicAuth() *Route {
	r.IsBasicAuth = true
	return r
}

// JWT Вешаем флаг авторизации JW Token
func (r *Route) JWT() *Route {
	r.IsBasicAuth = true
	return r
}

// Session Вешаем получение аутентификации сессии
func (r *Route) Session() *Route {
	r.IsSession = true
	return r
}

// Permission Вешаем получение аутентификации сессии
func (r *Route) Permission() *Route {
	r.IsPermission = true
	return r
}

// Login Вешаем получение аутентификации сессии
func (r *Route) Login(loginHandler AuthHandler, expires time.Time) *Route {
	r.Handler = func(ctx *fiber.Ctx) error {

		key := utils.UUID()
		err := loginHandler(ctx, key)
		if err != nil {
			return ctx.Status(401).SendString(err.Error())
		}

		cookie := new(fiber.Cookie)
		cookie.Name = sessionId
		cookie.Value = key
		cookie.Expires = expires
		ctx.Cookie(cookie)

		return ctx.SendStatus(200)
	}
	return r
}

// Logout Вешаем получение аутентификации сессии
func (r *Route) Logout(logoutHandler AuthHandler, route string) *Route {
	r.Handler = func(ctx *fiber.Ctx) error {

		key := ctx.Cookies(sessionId)
		err := logoutHandler(ctx, key)
		if err != nil {
			return ctx.Status(401).SendString(err.Error())
		}

		ctx.ClearCookie(sessionId)

		/*cookie := new(fiber.Cookie)
		cookie.Name = ""
		cookie.Expires = time.Now()
		ctx.Cookie(cookie)*/

		return ctx.Redirect(route)
	}
	return r
}
