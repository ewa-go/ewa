package auth

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type ApiKey struct {
	Handler      ApiKeyHandler
	Unauthorized ewa.EmptyHandler
}

func (a *ApiKey) Do(handler ewa.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ""
		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &ewa.Identity{
			User:   username,
			Domain: "",
			Permission: ewa.Permission{
				Route: ctx.Route(),
				//IsPermit: IsPermission,
			},
		})
	}
}
