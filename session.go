package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

// Session структура, которая описывает сессию
type Session struct {
	RedirectPath string
	AllRoutes    bool
	Expires      time.Duration
	Handler      SessionHandler
	ErrorHandler ErrorHandler
}

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	User   string
	Domain string
	//Permission Permission
}

/*type Permission struct {
	Route    *fiber.Route
	IsPermit bool
}*/

const sessionId = "session_id"

// Проверяем куки и извлекаем по ключу id по которому в бд/файле/памяти находим запись
func (s *Session) check(handler Handler, isPermission bool, permission *Permission) EmptyHandler {
	return func(ctx *fiber.Ctx) (err error) {

		user := "Unknown"
		domain := "Unknown"

		// Если cookie не существует, то перенаправляем запрос на условно "/login"
		key := ctx.Cookies(sessionId)
		if len(key) == 0 {
			return ctx.Redirect(s.RedirectPath)
		}

		// Проверяем на существование Handler
		if s.Handler != nil {
			user, domain, err = s.Handler(ctx.Cookies(sessionId))
			if err != nil {
				return ctx.Redirect(s.RedirectPath)
			}
		}

		// Получаем путь, чтобы передать в WebHandler
		route := ctx.Route()
		// Проверяем на существование PermissionHandler
		if isPermission && permission != nil && route != nil {
			if !permission.Handler(key, route.Path) {
				if s.ErrorHandler != nil {
					return s.ErrorHandler(ctx, fiber.StatusForbidden)
				}
				return ctx.SendStatus(fiber.StatusForbidden)
			}
		}

		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			User:   user,
			Domain: domain,
			/*Permission: Permission{
				Route:    route,
				IsPermit: IsPermission,
			},*/
		})
	}
}

// Формируем session_id и добавляем в куки
func (s *Session) login(handler WebAuthHandler) EmptyHandler {
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
func (s *Session) logout(handler WebAuthHandler) EmptyHandler {
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
