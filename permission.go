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

		route := ctx.Route()
		if !p.Check(route.Path) {
			return p.Error(ctx, 403, "Доступ запрещен (Permission denied)")
		}

		return handler(ctx)
	}
}
