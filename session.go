package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

// Session Структура которая описывает сессию
type Session struct {
	RedirectPath      string
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
