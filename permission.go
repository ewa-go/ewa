package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

// Permission Структура описывает разрешения на запрос
type Permission struct {
	Check PermissionHandler
	Error ErrorHandler
}

// Проверяем запрос на разрешения
func (p *Permission) check(handler Handler) Handler {
	return func(ctx *fiber.Ctx) error {

		key := ctx.Cookies(sessionId)
		route := ctx.Route()
		if !p.Check(key, route.Path) {
			if p.Error != nil {
				return p.Error(ctx, 403, "Доступ запрещен (Permission denied)")
			}
			return ctx.Status(403).SendString("Доступ запрещен (Permission denied)")
		}

		return handler(ctx)
	}
}
