package ewa

import (
	"encoding/json"
	"fmt"
	f "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/session"
	"github.com/gofiber/fiber/v2"
	"testing"
	"time"
)

func TestSession(t *testing.T) {

	cfg := Config{
		Port: 8877,
		Session: &session.Config{
			RedirectPath:         "/login",
			Expires:              24 * time.Hour,
			SessionHandler:       sessionHandler,
			DeleteSessionHandler: deleteSessionHandler,
		},
		Permission: nil,
		Static: &Static{
			Prefix: "/",
			Root:   "./views",
		},
		NotFoundPage: "",
		Views: &Views{
			Root:   "./views",
			Engine: f.Html,
		},
		ContextHandler: contextHandler,
		ErrorHandler:   nil,
	}
	app := fiber.New(fiber.Config{
		Views: f.NewViews("./views", f.Html, &f.Engine{
			Reload: true,
		}),
	})

	// Новый сервер
	server := New(&f.Server{App: app}, cfg)
	server.Register(new(Login)).SetPath("/login")
	server.Register(new(Home)).SetPath("/home")
	server.Register(new(Logout)).SetPath("/logout")

	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func contextHandler(handler Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(NewContext(&f.Context{Ctx: ctx}))
	}
}

func sessionHandler(value string) (string, error) {
	return "user", nil
}

func deleteSessionHandler(value string) bool {
	return true
}

type Home struct{}

func (Home) Get(route *Route) {
	route.Session().Handler = func(c *Context) error {
		return c.Render("home", nil)
	}
}

type Logout struct{}

func (Logout) Get(route *Route) {
	route.Session(Off).Handler = func(c *Context) error {
		return nil
	}
}

type Login struct{}

func (Login) Get(route *Route) {
	route.Handler = func(c *Context) error {
		return c.Render("login", nil)
	}
}

func (l Login) Post(route *Route) {
	route.Session(On).Handler = func(c *Context) error {

		body := c.Body()
		err := json.Unmarshal(body, &l)
		if err != nil {
			return c.SendString(400, err.Error())
		}

		var sessionId string
		if c.Session != nil {
			sessionId = c.Session.Value
		}

		fmt.Println(sessionId)

		return c.SendString(200, "OK")
	}
}

func TestNewSuffix(t *testing.T) {

	hostname := Suffix{
		Index:       2,
		Value:       "hostname",
		isParam:     false,
		Description: "Hostname",
	}

	client := Suffix{
		Index:       2,
		Value:       "client",
		isParam:     false,
		Description: "Client",
	}

	s := NewSuffix(hostname, client)
	fmt.Printf("%#v", s)
}
