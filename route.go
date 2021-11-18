package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

type Route struct {
	Params        []string
	Authorization int
	Handler       interface{}
	IsSession     bool
	IsPermission  bool
	webAuth       *WebAuth
	webSocket     *WebSocket
	Option        Option
}

type WebSocket struct {
	UpgradeHandler fiber.Handler
	Handler        WsHandler
}

type Option struct {
	Name        string
	Description string
	Body        string
}

type WebAuth struct {
	IsLogin bool
	Handler WebAuthHandler
}

const (
	NoAuth = iota
	BasicAuth
	DigestAuth
	ApiKeyAuth
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
func (r *Route) Auth(auth int) *Route {
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
func (r *Route) WebSocket(upgrade fiber.Handler) *WebSocket {
	r.webSocket = &WebSocket{
		UpgradeHandler: upgrade,
	}
	return r.webSocket
}

// WebAuth Устанавливаем web socket соединение
func (r *Route) WebAuth(isLogin bool) *WebAuth {
	r.webAuth = &WebAuth{
		IsLogin: isLogin,
	}
	return r.webAuth
}

func (r *Route) Empty() {
	r.Handler = nil
}

// SetOption устанавливаем опции для свагера
func (r *Route) SetOption(name, description, body string) *Route {
	r.Option = Option{
		Name:        name,
		Description: description,
		Body:        body,
	}
	return r
}

func (r *Route) GetHandler(config Config) fiber.Handler {
	switch r.Handler.(type) {
	// handler для маршрутов с identity
	case Handler:
		// Авторизация
		switch r.Authorization {
		case NoAuth:
			return func(ctx *fiber.Ctx) error {
				return r.Handler.(Handler)(ctx, nil)
			}
		case BasicAuth:
			if config.Authorization.Basic != nil {
				return config.Authorization.Basic.Do(r.Handler.(Handler))
			}
			break
		case DigestAuth:
			if config.Authorization.Digest != nil {
				return config.Authorization.Digest.Do(r.Handler.(Handler))
			}
			break
		case ApiKeyAuth:
			if config.Authorization.ApiKey != nil {
				return config.Authorization.ApiKey.Do(r.Handler.(Handler))
			}
			break
		}
		break

	// handler для маршрутов web авторизации Login и Logout
	case WebAuthHandler:
		if config.Session != nil && r.webAuth != nil {
			if r.webAuth.IsLogin {
				// Авторизация - вход
				return config.Session.login(r.Handler.(WebAuthHandler))
			}
			// Авторизация - выход
			return config.Session.logout(r.Handler.(WebAuthHandler))
		}
	}

	return nil
}
