package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	RedirectPath string
	Check        CheckHandler
}

func (s *Session) check(handler Handler) Handler {
	return func(ctx *fiber.Ctx) error {

		key := ctx.Cookies(sessionId)
		user, err := s.Check(key)
		if err != nil {
			return ctx.Redirect(s.RedirectPath)
		}

		ctx.Set("User", user)

		return handler(ctx)
	}
}
