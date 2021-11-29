package egowebapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Route struct {
	Params        []string    `json:"params,omitempty"`
	Authorization Auth        `json:"authorization"`
	Handler       interface{} `json:"-"`
	IsSession     bool        `json:"is_session"`
	IsPermission  bool        `json:"is_permission"`
	WebAuth       *WebAuth    `json:"web_auth,omitempty"`
	webSocket     *WebSocket
	Option        Option `json:"option"`
}

type WebSocket struct {
	UpgradeHandler fiber.Handler
}

type WebAuth struct {
	IsLogin bool `json:"is_login"`
}

type Option struct {
	Headers     []string `json:"headers,omitempty"`
	Method      string   `json:"method,omitempty"`
	Body        string   `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Auth string

const (
	NoAuth     Auth = "NoAuth"
	BasicAuth  Auth = "BasicAuth"
	DigestAuth Auth = "DigestAuth"
	ApiKeyAuth Auth = "ApiKeyAuth"
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
func (r *Route) Auth(auth Auth) *Route {
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

// GetHandler возвращаем обработчик основанный на параметрах конфигурации маршрута
func (r *Route) GetHandler(s *Server) fiber.Handler {

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
			if s.Config.Authorization.Basic != nil {
				return s.Config.Authorization.Basic.Do(h, r.IsPermission, s.Config.Permission)
			}
			break
		case DigestAuth:
			if s.Config.Authorization.Digest != nil {
				return s.Config.Authorization.Digest.Do(h, r.IsPermission, s.Config.Permission)
			}
			break
		case ApiKeyAuth:
			if s.Config.Authorization.ApiKey != nil {
				return s.Config.Authorization.ApiKey.Do(h, r.IsPermission, s.Config.Permission)
			}
			break
		}

		// Проверяем маршрут на актуальность сессии
		if (r.IsSession && s.Config.Session != nil) || r.IsSession {
			return s.Config.Session.check(h, r.IsPermission, s.Config.Permission)
		}
		return func(ctx *fiber.Ctx) error {
			return h(ctx, nil)
		}

	// Swagger handler для добавления описания маршрутов
	case func(*fiber.Ctx, *Swagger) error:
		return func(ctx *fiber.Ctx) error {
			return h(ctx, s.Swagger)
		}

	// Handler для маршрутов web авторизации Login и Logout
	case func(*fiber.Ctx, string) error:
		if s.Config.Session != nil && r.WebAuth != nil {
			if r.WebAuth.IsLogin {
				// Авторизация - вход
				return s.Config.Session.login(h)
			}
			// Авторизация - выход
			return s.Config.Session.logout(h)
		}
		break

	// Handler для маршрут WebSocket соединения
	case func(*websocket.Conn):
		return websocket.New(h)

	// Обычный обработчик без ништяков
	case func(*fiber.Ctx) error:
		return h
	}

	// Ну если ни один из обработчиков не удовлетворяет требованиям, то вернем ответ с кодом 404
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("%s %s", ctx.Route().Method, ctx.Route().Path))
	}
}
