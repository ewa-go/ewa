package egowebapi

type Handler func(c *Context) error
type PermissionHandler func(username string, path string) bool
type SessionHandler func(key string) (user string, err error)
type GenSessionIdHandler func() string
type ErrorHandler func(c *Context, statusCode int, err interface{}) error
type SignHandler func(c *Context, key string) error

//type SwaggerHandler func(c *Context, swagger *s.Swagger) error

// Авторизация
