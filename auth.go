package egowebapi

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strings"
)

const (
	Md5Algorithm           = "MD5"
	Md5SessAlgorithm       = "MD5-sess"
	Sha256Algorithm        = "SHA-256"
	Sha256SessAlgorithm    = "SHA-256-sess"
	Sha512256Algorithm     = "SHA-512-256"
	Sha512256SessAlgorithm = "SHA-512-256-sess"
)

//const StatusForbidden = "Доступ запрещен (Permission denied)"

type Authorization struct {
	AllRoutes Auth
	Basic     *Basic
	Digest    *Digest
	ApiKey    *ApiKey
}

type Basic struct {
	Handler      BasicAuthHandler
	Unauthorized ErrorHandler
}

func (b *Basic) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	i := strings.IndexByte(cs, ':')
	if i < 0 {
		return
	}
	return cs[:i], cs[i+1:], true
}

func (b *Basic) realm(ctx *fiber.Ctx) error {
	if b.Unauthorized == nil {
		ctx.Set(fiber.HeaderWWWAuthenticate, `Basic realm="Необходимо указать имя пользователя и пароль"`)
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	return b.Unauthorized(ctx, fiber.StatusUnauthorized)
}

func (b *Basic) Do(handler Handler, isPermission bool, permission *Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get(fiber.HeaderAuthorization)
		if auth == "" {
			return b.realm(ctx)
		}

		username, password, ok := b.parseBasicAuth(auth)
		if !ok || !b.Handler(username, password) {
			return b.realm(ctx)
		}

		// Получаем путь
		route := ctx.Route()
		// Проверяем на существование PermissionHandler
		if isPermission && permission != nil && route != nil {
			if !permission.Handler(username, route.Path) {
				if b.Unauthorized != nil {
					return b.Unauthorized(ctx, fiber.StatusForbidden)
				}
				return ctx.SendStatus(fiber.StatusForbidden)
			}
		}

		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			User:   username,
			Domain: "",
			/*Permission: Permission{
				Route: ctx.Route(),
				//IsPermit: IsPermission,
			},*/
		})
	}
}

type Digest struct {
	Handler      DigestAuthHandler
	Unauthorized ErrorHandler
}

type Advanced struct {
	Realm       string
	Nonce       string
	Algorithm   string
	Qop         string
	NonceCount  string
	ClientNonce string
	Opaque      string
}

func (d *Digest) Do(handler Handler, isPermission bool, permission *Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ""
		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			User:   username,
			Domain: "",
			/*Permission: Permission{
				Route: ctx.Route(),
				//IsPermit: IsPermission,
			},*/
		})
	}
}

type ApiKey struct {
	Handler      ApiKeyHandler
	Unauthorized ErrorHandler
}

func (a *ApiKey) Do(handler Handler, isPermission bool, permission *Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ""
		// Возвращаем данные по пользователю и маршруту
		return handler(ctx, &Identity{
			User:   username,
			Domain: "",
			/*Permission: Permission{
				Route: ctx.Route(),
				//IsPermit: IsPermission,
			},*/
		})
	}
}
