package egowebapi

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type BasicAuth struct {
	Authorizer   Authorizer
	Unauthorized Handler
}

func NewBasicAuth(authorizer Authorizer, unauthorized Handler) *BasicAuth {
	return &BasicAuth{
		Authorizer:   authorizer,
		Unauthorized: unauthorized,
	}
}

func (b *BasicAuth) parseBasicAuth(auth string) (username, password string, ok bool) {
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

func (b *BasicAuth) check(handler Handler) Handler {
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
		if !ok || !b.Authorizer(username, password) {
			if b.Unauthorized == nil {
				ctx.Set("WWW-Authenticate", `Basic realm="Необходимо указать имя пользователя и пароль"`)
				return ctx.SendStatus(fiber.StatusUnauthorized)
			}
			return b.Unauthorized(ctx)
		}

		return handler(ctx)
	}
}
