package egowebapi

import "github.com/gofiber/fiber/v2"

/*
type IController interface {
	IHttp
	IRest
}*/

type IWeb interface {
	Get() *Route
	Post() *Route
}

type IRest interface {
	IWeb
	Put() *Route
	Delete() *Route
	Options(swagger *Swagger) Handler
}

type Handler fiber.Handler

type Route struct {
	Params      []string
	Description string
	Handler     Handler
}

type Options []*Option

type Option struct {
	Params      []string
	Description string
	Method      string
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
	s.Options = append(s.Options, option)
}

func SetParams(params ...string) []string {
	return params
}

func NewRoute(handler Handler, params ...string) *Route {
	return &Route{
		Params:  params,
		Handler: handler,
	}
}

func (r *Route) SetDescription(s string) *Route {
	r.Description = s
	return r
}
