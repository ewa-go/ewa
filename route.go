package egowebapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Route struct {
	Params        []string    `json:"-"`
	Authorization string      `json:"authorization"`
	Handler       interface{} `json:"-"`
	IsSession     bool        `json:"is_session"`
	IsPermission  bool        `json:"is_permission"`
	WebAuth       *WebAuth    `json:"web_auth,omitempty"`
	Option        Option      `json:"option"`
	webSocket     *WebSocket
}

type WebSocket struct {
	UpgradeHandler fiber.Handler
}

type WebAuth struct {
	IsLogin bool `json:"is_login"`
}

type Option struct {
	Headers     []string `json:"headers"`
	Method      string   `json:"method"`
	Body        string   `json:"body"`
	Description string   `json:"description"`
}

const (
	NoAuth     = "NoAuth"
	BasicAuth  = "BasicAuth"
	DigestAuth = "DigestAuth"
	ApiKeyAuth = "ApiKeyAuth"
)

// SetParams указываем параметры маршрута
func (r *Route) SetParams(params ...string) *Route {
	r.Params = params
	return r
}

// SetDescription устанавливаем описание маршрута
func (r *Route) SetDescription(s string) *Route {
	r.Option.Description = s
	return r
}

// SetBody устанавливаем описание тела маршрута
func (r *Route) SetBody(s string) *Route {
	r.Option.Body = s
	return r
}

// Auth указываем метод авторизации
func (r *Route) Auth(auth string) *Route {
	r.Authorization = auth
	return r
}

// Session вешаем получение аутентификации сессии,
func (r *Route) Session() *Route {
	r.IsSession = true
	return r
}

// Permission ставим флаг для проверки маршрута на право доступа
func (r *Route) Permission() *Route {
	r.IsPermission = true
	return r
}

// WebSocket Устанавливаем web socket соединение
func (r *Route) WebSocket(upgrade fiber.Handler) *Route {
	r.webSocket = &WebSocket{
		UpgradeHandler: upgrade,
	}
	return r
}

// SetWebAuth Устанавливаем web socket соединение
func (r *Route) SetWebAuth(isLogin bool) *Route {
	r.WebAuth = &WebAuth{
		IsLogin: isLogin,
	}
	return r
}

func (r *Route) Empty() {
	r.Handler = nil
}

// SetOption устанавливаем опции для свагера
func (r *Route) SetOption(name, description, body string) *Route {
	r.Option = Option{
		//Name:        name,
		Description: description,
		Body:        body,
	}
	return r
}

func (r *Route) GetHandler(config Config) fiber.Handler {

	switch h := r.Handler.(type) {
	// handler для маршрутов с identity
	case func(*fiber.Ctx, *Identity) error:
		// Авторизация
		switch r.Authorization {
		case NoAuth:
			return func(ctx *fiber.Ctx) error {
				return h(ctx, nil)
			}
		case BasicAuth:
			if config.Authorization.Basic != nil {
				return config.Authorization.Basic.Do(h, r.IsPermission, config.Permission)
			}
			break
		case DigestAuth:
			if config.Authorization.Digest != nil {
				return config.Authorization.Digest.Do(h, r.IsPermission, config.Permission)
			}
			break
		case ApiKeyAuth:
			if config.Authorization.ApiKey != nil {
				return config.Authorization.ApiKey.Do(h, r.IsPermission, config.Permission)
			}
			break
		}

		// Проверяем маршрут на актуальность сессии
		if (r.IsSession && config.Session != nil) || r.IsSession {
			return config.Session.check(h, r.IsPermission, config.Permission)
		}
		return func(ctx *fiber.Ctx) error {
			return h(ctx, nil)
		}

	// Swagger handler для добавления описания маршрутов
	case SwaggerHandler:

		h = s.Swagger.check(r.Handler.(SwaggerHandler))

		break
		// handler для маршрутов web авторизации Login и Logout
	case func(*fiber.Ctx, string) error:
		if config.Session != nil && r.WebAuth != nil {
			if r.WebAuth.IsLogin {
				// Авторизация - вход
				return config.Session.login(h)
			}
			// Авторизация - выход
			return config.Session.logout(h)
		}
		break

	// Handler для маршрут WebSocket соединения
	case func(*websocket.Conn):
		if r.webSocket != nil {
			if r.webSocket.UpgradeHandler != nil {
				return r.webSocket.UpgradeHandler
			}
			return websocket.New(h)
		}
		break
	}

	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("%s %s", ctx.Route().Method, ctx.Route().Path))
	}
}
