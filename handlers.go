package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type PermissionHandler func(key, path string) bool
type SessionHandler func(key string) (user string, domain string, err error)
type ErrorHandler func(ctx *fiber.Ctx, code int, err string) error
type AuthHandler func(ctx *fiber.Ctx, key string) error
type Handler fiber.Handler
type WebHandler func(ctx *fiber.Ctx, identity *Identity) error
type WsHandler func(c *websocket.Conn)
type Authorizer func(user string, pass string) bool
