package auth

import (
	"encoding/base64"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Basic struct {
	Handler      BasicAuthHandler
	Unauthorized ewa.EmptyHandler
}

func (b *Basic) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	i := strings.IndexByte(cs, ':')
	if i < 0 {
		return
	}
	return cs[:i], cs[i+1:], true
}

func (b *Basic) Do(handler ewa.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if auth == "" {
			if b.Unauthorized == nil {
				ctx.Set("WWW-Authenticate", `Basic realm="Необходимо указать имя пользователя и пароль"`)
				return ctx.SendStatus(fiber.StatusUnauthorized)
			}
			return b.Unauthorized(ctx)
		}

		username, password, ok := b.parseBasicAuth(auth)
		if !ok || !b.Handler(username, password) {
			if b.Unauthorized == nil {
				ctx.Set("WWW-Authenticate", `Basic realm="Необходимо указать имя пользователя и пароль"`)
				return ctx.SendStatus(fiber.StatusUnauthorized)
			}
			return b.Unauthorized(ctx)
		}

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
