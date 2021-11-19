package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type EmptyHandler func(ctx *fiber.Ctx) error
type PermissionHandler func(id interface{}, path string) bool
type SessionHandler func(key string) (user string, domain string, err error)
type ErrorHandler func(ctx *fiber.Ctx, statusCode int) error
type WebAuthHandler func(ctx *fiber.Ctx, key string) error
type Handler func(ctx *fiber.Ctx, identity *Identity) error
type WsHandler func(c *websocket.Conn)
type SwaggerHandler func(ctx *fiber.Ctx, swagger *Swagger) error
type BasicAuthHandler func(user string, pass string) bool
type DigestAuthHandler func(user string, pass string, advanced Advanced) bool
type ApiKeyHandler func(key string, value string) bool
