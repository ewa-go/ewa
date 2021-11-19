package egowebapi

type IGet interface {
	Get(route *Route)
}

type IPost interface {
	Post(route *Route)
}
type IPut interface {
	Put(route *Route)
}

type IDelete interface {
	Delete(route *Route)
}

type IWeb interface {
	IGet
	IPost
}

type IRest interface {
	IWeb
	IPut
	IDelete
}

type IRestOptions interface {
	IRest
	Options(route *Route)
}
