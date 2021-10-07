package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

// Session Структура которая описывает сессию
type Session struct {
	RedirectPath      string
	Expires           time.Duration
	SessionHandler    SessionHandler
	PermissionHandler PermissionHandler
	ErrorHandler      ErrorHandler
}

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	User       string
	Domain     string
	Permission Permission
}

type Permission struct {
	Route    *fiber.Route
	IsPermit bool
}

const StatusForbidden = "Доступ запрещен (Permission denied)"
const sessionId = "session_id"

// Проверяем куки и извлекаем по ключу id по которому в бд находим запись
func (s *Session) check(handler WebHandler, IsPermission bool) Handler {
	return func(ctx *fiber.Ctx) (err error) {

		user := "Unknown"
		domain := "Unknown"

		// Если cookie не существует, то перенаправляем запрос на условно "/login"
		key := ctx.Cookies(sessionId)
		if len(key) == 0 {
			return ctx.Redirect(s.RedirectPath)
		}

		// Получаем путь, чтобы передать в WebHandler
		route := ctx.Route()

		// Проверяем на существование SessionHandler
		if s.SessionHandler != nil {
			user, domain, err = s.SessionHandler(ctx.Cookies(sessionId))
			if err != nil {
				return ctx.Redirect(s.RedirectPath)
			}
		}

		// Проверяем на существование PermissionHandler
		if IsPermission && s.PermissionHandler != nil {
			if !s.PermissionHandler(key, route.Path) {
				if s.ErrorHandler != nil {
					return s.ErrorHandler(ctx, 403, StatusForbidden)
				}
				return ctx.Status(403).SendString(StatusForbidden)
			}
		}

		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			User:   user,
			Domain: domain,
			Permission: Permission{
				Route:    route,
				IsPermit: IsPermission,
			},
		})
	}
}

// Формируем session_id и добавляем в куки
func (s *Session) login(handler AuthHandler) Handler {
	return func(ctx *fiber.Ctx) error {

		key := utils.UUID()
		err := handler(ctx, key)
		if err != nil {
			return ctx.Status(401).SendString(err.Error())
		}

		cookie := new(fiber.Cookie)
		cookie.Name = sessionId
		cookie.Value = key
		cookie.Expires = time.Now().Add(s.Expires)
		ctx.Cookie(cookie)

		return ctx.SendStatus(200)
	}
}

// Очищаем куки, чтобы при маршрутизации сессия не была доступна
func (s *Session) logout(handler AuthHandler) Handler {
	return func(ctx *fiber.Ctx) error {

		key := ctx.Cookies(sessionId)
		err := handler(ctx, key)
		if err != nil {
			return ctx.Status(401).SendString(err.Error())
		}

		ctx.ClearCookie(sessionId)

		return ctx.Redirect(s.RedirectPath)
	}

}
