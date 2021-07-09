package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

type Permission struct {
	Check PermissionHandler
	Error ErrorHandler
}

func (p *Permission) check(handler Handler) Handler {
	return func(ctx *fiber.Ctx) error {

		key := ctx.Cookies(sessionId)
		route := ctx.Route()
		if !p.Check(key, route.Path) {
			return p.Error(ctx, 403, "Доступ запрещен (Permission denied)")
		}

		return handler(ctx)
	}
}
