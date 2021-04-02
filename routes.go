package egowebapi

import "github.com/gofiber/fiber/v2"

type Handler fiber.Handler

type Route struct {
	Path    []string
	Handler Handler
}

func AddPath(path ...string) []string {
	return path
}

func NewRoute(handler Handler, path ...string) *Route {
	return &Route{
		Path:    path,
		Handler: handler,
	}
}
