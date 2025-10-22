package ewa

import (
	"fmt"
	"testing"
	"time"

	f "github.com/ewa-go/ewa-fiber"
	"github.com/ewa-go/ewa/v2/consts"
	"github.com/ewa-go/jsonschema"
	"github.com/gofiber/fiber/v2"
)

type Test struct{}

func (Test) Get(route *Route) {
	route.SetParameters(NewPathParam("/{id}")).SetEmptyParam("get data")
	route.Handler = func(c *Context) error {

		id := c.Params("id")
		if len(id) > 0 {
			return c.JSON(200, id)
		}
		return c.JSON(200, []string{id})
	}
}

type API struct{}

func (API) Get(route *Route) {
	route.Handler = func(c *Context) error {
		b, err := c.Swagger.JSON()
		if err != nil {
			return c.SendString(consts.StatusBadRequest, err.Error())
		}
		return c.Send(consts.StatusOK, consts.MIMEApplicationJSON, b)
	}
}

func newSwagger() *Swagger {
	s := &Swagger{
		Swagger:             "2.0",
		Host:                fmt.Sprintf("localhost:%d", 8877),
		BasePath:            "",
		SecurityDefinitions: SecurityDefinitions{},
		Paths:               Paths{},
		Definitions:         jsonschema.Definitions{},
		models:              Models{},
	}
	s.SetBasePath("").SetInfo("localhost", &Info{
		Description: "Description",
		Version:     "0.0.1",
		Title:       "Title",
		Contact: &Contact{
			Email: "yegor.govorukhin@mail.ru",
		},
		License: &License{
			Name: "Freeware license",
		},
	}, nil)
	return s
}

func newServer() *Server {

	cfg := Config{
		Port: 8877,
		Session: &Session{
			RedirectPath:         "/login",
			Expires:              24 * time.Hour,
			SessionHandler:       sessionHandler,
			DeleteSessionHandler: deleteSessionHandler,
		},
		Permission: &Permission{
			AllRoutes: true,
			Handler: func(c *Context, identity *Identity, method, path string) bool {
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
		Authorization: Authorization{
			Unauthorized: func(err error) bool {
				fmt.Println(err)
				return true
			},
			Basic: &Basic{
				Handler: func(c *Context, user string, pass string) error {
					if user == "user" && pass == "Qq123456" {
						return nil
					}
					return nil
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
	server := New(&f.Server{App: app}, cfg)
	server.Register(new(API))
	server.Swagger = newSwagger()

	return server
}

func contextHandler(handler Handler) interface{} {
	return func(ctx *fiber.Ctx) error {
		return handler(NewContext(f.IContext(ctx)))
	}
}

func sessionHandler(c *Context, value string) (string, error) {
	fmt.Println("sessionHandler", value)
	return "user", nil
}

func deleteSessionHandler(value string) bool {
	fmt.Println("deleteSessionHandler", value)
	return true
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

func TestServer(t *testing.T) {

	s := newServer()
	schema := &Suffix{
		Index: 1,
		Value: "schema",
	}
	user := &Suffix{
		Index: 2,
		Value: "{user}name",
	}
	s.Register(new(Test)).SetSuffix(schema, user)
	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}
}
