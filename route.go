package egowebapi

import "github.com/gofiber/fiber/v2"

type Route struct {
	Params      []string
	Description string
	IsBasicAuth bool
	Handler     Handler
}

type Handler fiber.Handler

func (r *Route) SetHandler(handler Handler) *Route {
	r.Handler = handler
	return r
}

func (r *Route) SetParams(params ...string) *Route {
	r.Params = params
	return r
}

func (r *Route) SetDescription(s string) *Route {
	r.Description = s
	return r
}

func (r *Route) SetBasicAuth(b bool) *Route {
	r.IsBasicAuth = b
	return r
}
