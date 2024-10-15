package ewa

import (
	"errors"
	"fmt"
	f "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/consts"
	"github.com/ewa-go/ewa/security"
	"github.com/ewa-go/ewa/session"
	"github.com/gofiber/fiber/v2"
	"testing"
	"time"
)

type About struct{}

func (About) Get(route *Route) {
	route.SetSecurity(security.BasicAuth, security.ApiKeyAuth)
	route.Handler = func(c *Context) error {
		return c.SendString(200, "О программе")
	}
}

func startServer() error {

	cfg := Config{
		Port: 8877,
		Session: &session.Config{
			RedirectPath:         "/login",
			Expires:              24 * time.Hour,
			SessionHandler:       sessionHandler,
			DeleteSessionHandler: deleteSessionHandler,
		},
		Permission: &Permission{
			AllRoutes: true,
			Handler: func(c *Context, identity *security.Identity, method, path string) bool {
				if identity != nil && identity.Username == "user" {
					// Set admin variable
					identity.SetVariable("is_admin", false)
					switch method {
					// ReadOnly
					case consts.MethodPost, consts.MethodPut, consts.MethodDelete:
						return false
					}
				}
				return true
			},
			NotPermissionHandler: nil,
		},
		NotFoundPage:   "",
		ContextHandler: contextHandler,
		ErrorHandler:   nil,
		Authorization: security.Authorization{
			Unauthorized: func(err error) bool {
				fmt.Println(err)
				return true
			},
			Basic: &security.Basic{
				Handler: func(user string, pass string) error {
					if user == "user" && pass == "Qq123456" {
						return nil
					}
					return errors.New("username or password not correct")
				},
			},
			ApiKey: &security.ApiKey{
				KeyName: "Token",
				Param:   security.ParamHeader,
				Handler: func(token string) (username string, err error) {
					if token == "cb96a323-6d6b-44ce-9f40-5c3f9e365800" {
						return "user", nil
					}
					return "", errors.New("token is invalid")
				},
			},
		},
	}
	app := fiber.New()

	// Новый сервер
	server := New(&f.Server{App: app}, cfg)

	server.Register(new(About)).SetPath("/about")

	return server.Start()
}

func contextHandler(handler Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(NewContext(f.IContext(ctx)))
	}
}

func sessionHandler(value string) (string, error) {
	fmt.Println("sessionHandler", value)
	return "user", nil
}

func deleteSessionHandler(value string) bool {
	fmt.Println("deleteSessionHandler", value)
	return true
}

func TestServer(t *testing.T) {
	err := startServer()
	if err != nil {
		t.Fatal(err)
	}
}
