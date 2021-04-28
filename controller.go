package egowebapi

import "github.com/gofiber/fiber/v2"

/*
type IController interface {
	IHttp
	IRest
}*/

type IWeb interface {
	Get(route *Route)
	Post(route *Route)
}

type IRest interface {
	IWeb
	Put(route *Route)
	Delete(route *Route)
	Options(swagger *Swagger) Handler
}

type IController interface {
	Path() string
}

type Controller struct {
	Path string
}

//Swagger

type Options []*Option

type Option struct {
	Params      []string
	Headers     []string
	Method      string
	Body        string
	Description string
}

type Swagger struct {
	Name    string
	Path    string
	Options Options
}

func newSwagger(name string, path string) *Swagger {
	return &Swagger{
		Name: name,
		Path: path,
	}
}

func (s *Swagger) AddOption(option *Option) {
	if option != nil {
		s.Options = append(s.Options, option)
	}
}

func (s *Swagger) Allow(ctx *fiber.Ctx) {
	var methods []string
	for _, option := range s.Options {
		methods = append(methods, option.Method)
	}
	ctx.Append("Allow", methods...)
}
