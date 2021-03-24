package egowebapi

import (
	"github.com/gofiber/fiber"
)

type Route struct {
	Method      string
	Path        string
	Description []string
	Handler     fiber.Handler
}

type Routes []*Route

func NewRoute(method string, path string, handler fiber.Handler, description ...string) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Description: description,
	}
}
