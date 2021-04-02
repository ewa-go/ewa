package egowebapi

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
	Options() *Route
}
