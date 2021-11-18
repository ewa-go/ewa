package egowebapi

import (
	"github.com/gofiber/fiber/v2"
)

type Swagger struct {
	Uri    string     `json:"uri"`
	Routes []RouteExt `json:"routes"`
}

type RouteExt struct {
	Name string `json:"name"`
	Path string `json:"path"`
	*Route
}

func (s *Swagger) Add(path string, route *Route) {
	s.Routes = append(s.Routes, RouteExt{
		Path:  s.Uri + path,
		Route: route,
	})
}

func (s *Swagger) Allow(ctx *fiber.Ctx) {
	/*var methods []string
	for _, option := range s.Options {
		methods = append(methods, option.Method)
	}
	ctx.Append("Allow", methods...)*/
}

/*func (s *Swagger) check(handler SwaggerHandler) Handler {
	return func(ctx *fiber.Ctx) error {
		return handler(ctx, s)
	}
}*/
