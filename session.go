package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

// Session Структура которая описывает сессию
type Session struct {
	RedirectPath string
	Check        CheckHandler
}

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	User   string
	Domain string
}

// Проверяем куки и извлекаем по ключу id по которому в бд находим запись
func (s *Session) check(handler WebHandler) Handler {
	return func(ctx *fiber.Ctx) error {

		user, err := s.Check(ctx.Cookies(sessionId))
		if err != nil {
			return ctx.Redirect(s.RedirectPath)
		}

		//ctx.Set("User", user)

		return handler(ctx, &Identity{User: user})
	}
}
