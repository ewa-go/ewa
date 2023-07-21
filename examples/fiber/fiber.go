package main

import (
	"fiber/controllers/api"
	"fiber/controllers/web"
	"fmt"
	"github.com/ewa-go/ewa"
	f "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/security"
	"github.com/ewa-go/ewa/session"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func main() {

	cfg := ewa.Config{
		Port: 8877,
		Session: &session.Config{
			RedirectPath:         "/login",
			Expires:              24 * time.Hour,
			SessionHandler:       sessionHandler,
			DeleteSessionHandler: deleteSessionHandler,
		},
		Permission: nil,
		Static: &ewa.Static{
			Prefix: "/",
			Root:   "./views",
		},
		NotFoundPage: "",
		Views: &ewa.Views{
			Root:   "./views",
			Engine: f.Html,
		},
		ContextHandler: contextHandler,
		ErrorHandler:   nil,
		Authorization: security.Authorization{
			Unauthorized: func(err error) bool {
				fmt.Println(err)
				return true
			},
			Basic: &security.Basic{
				Handler: func(user string, pass string) bool {
					if user == "user" && pass == "Qq123456" {
						return true
					}
					return false
				},
			},
			Digest: nil,
			ApiKey: nil,
			OAuth2: nil,
		},
	}
	app := fiber.New(fiber.Config{
		Views: f.NewViews("./views", f.Html, &f.Engine{
			Reload: true,
		}),
	})

	// Новый сервер
	server := ewa.New(&f.Server{App: app}, cfg)
	// Страницы
	server.Register(new(web.Login)).SetPath("/login").NotShow()
	server.Register(new(web.Home)).SetPath("/home").NotShow()
	server.Register(new(web.Logout)).SetPath("/logout").NotShow()
	//api
	server.Register(new(api.User))

	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func contextHandler(handler ewa.Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(ewa.NewContext(f.IContext(ctx)))
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
