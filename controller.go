package egowebapi

/*
type IController interface {
	IHttp
	IRest
}*/

type IGet interface {
	Get(route *Route)
}

type IWeb interface {
	IGet
	Post(route *Route)
}

type IRest interface {
	IWeb
	Put(route *Route)
	Delete(route *Route)
	Options(swagger *Swagger) EmptyHandler
}
