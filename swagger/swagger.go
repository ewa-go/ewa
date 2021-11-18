package swagger

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Swagger struct {
	Routes []*Route `json:"routes"`
}

type Route struct {
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Options     Options `json:"options"`
}

type Option struct {
	Headers     []string `json:"headers"`
	Method      string   `json:"method"`
	Body        string   `json:"body"`
	Description string   `json:"description"`
}

type Options []*Option

func newSwagger(name string, path string) *Swagger {
	return &Swagger{
		//Name: name,
		//Path: path,
	}
}

func (s *Swagger) AddOption(option *Option) {
	if option != nil {
		//s.Options = append(s.Options, option)
	}
}

func (s *Swagger) Allow(ctx *fiber.Ctx) {
	/*var methods []string
	for _, option := range s.Options {
		methods = append(methods, option.Method)
	}
	ctx.Append("Allow", methods...)*/
}

func (s *Swagger) check(handler ewa.SwaggerHandler) ewa.Handler {
	return func(ctx *fiber.Ctx) error {
		return handler(ctx, s)
	}
}
