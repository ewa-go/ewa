package egowebapi

import "github.com/gofiber/fiber/v2"

type PermissionHandler func(key, path string) bool
type CheckHandler func(key string) (string, error)
type ErrorHandler func(ctx *fiber.Ctx, code int, err string) error
type AuthHandler func(ctx *fiber.Ctx, key string) error
type Handler fiber.Handler
type WebHandler func(ctx *fiber.Ctx, identity *Identity) error
type Authorizer func(user string, pass string) bool
