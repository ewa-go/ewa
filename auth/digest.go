package auth

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Digest struct {
	Handler      ewa.DigestAuthHandler
	Unauthorized ewa.EmptyHandler
}

type Advanced struct {
	Realm       string
	Nonce       string
	Algorithm   string
	Qop         string
	NonceCount  string
	ClientNonce string
	Opaque      string
}

func (d *Digest) Do(handler ewa.Handler) fiber.Handler {
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
