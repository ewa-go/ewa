package egowebapi

import "github.com/valyala/fasthttp"

const Empty = "У данного маршрута нет реализации"

type Method string

const (
	GET     Method = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
)

type Route struct {
	Path string
	Method string
	Description string
	Handler fasthttp.RequestHandler
}

type Routes []*Route

func NewRoutes(route ...*Route) (routes Routes) {
	routes = append(routes, route...)
	return routes
}

func NewRoute(path string, method string, description string, handler fasthttp.RequestHandler) *Route {
	if path == "" {
		path = "/"
	}
	return &Route{
		Path:        path,
		Method:      method,
		Description: description,
		Handler:     handler,
	}
}
