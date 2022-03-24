package egowebapi

type Handler func(c *Context) error
type PermissionHandler func(username string, path string) bool
type SessionHandler func(key string) (user string, err error)
type GenSessionIdHandler func() string
type ErrorHandler func(c *Context, statusCode int, err interface{}) error
type SignHandler func(c *Context, key string) error
type SwaggerHandler func(c *Context, swagger *Swagger) error

// Авторизация

type BasicAuthHandler func(user string, pass string) bool
type DigestAuthHandler func(user string, pass string, advanced Advanced) bool
type ApiKeyAuthHandler func(token string) (username string, err error)
