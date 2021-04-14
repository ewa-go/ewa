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

type Handler fiber.Handler

//Route

type Route struct {
	Params      []string
	Description string
	Handler     Handler
}

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

//Swagger

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
