package egowebapi

type Handler func(c *Context) error
type PermissionHandler func(id interface{}, path string) bool
type SessionHandler func(key string) (user string, domain string, err error)
type GenSessionIdHandler func() string
type ErrorHandler func(c *Context, statusCode int, err interface{}) error
type SignHandler func(c *Context, key string) error

//type WsHandler func(c *websocket.Conn)
type SwaggerHandler func(c *Context, swagger *Swagger) error
type BasicAuthHandler func(user string, pass string) bool
type DigestAuthHandler func(user string, pass string, advanced Advanced) bool
type ApiKeyAuthHandler func(key string, value string) bool
