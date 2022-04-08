package egowebapi

import v2 "github.com/egovorukhin/egowebapi/swagger/v2"

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

type IOptions interface {
	Options(route *Route)
}

type IPatch interface {
	Patch(route *Route)
}

type IHead interface {
	Head(route *Route)
}

type ITrace interface {
	Trace(route *Route)
}

type IConnect interface {
	Connect(route *Route)
}

type ITag interface {
	Tag() v2.Tag
}

/*type IWeb interface {
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
}*/
