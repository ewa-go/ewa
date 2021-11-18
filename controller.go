package egowebapi

import "github.com/egovorukhin/egowebapi/swagger"

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
	Options(swagger *swagger.Swagger) Handler
}

type IWebSocket interface {
	Get(route *Route)
}
