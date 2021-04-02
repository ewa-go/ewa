package api

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber"
)

var users = map[string]*User{}

type User struct {
	Lastname  string
	Firstname string
}

func (u *User) Get() *egowebapi.Route {
	return egowebapi.NewRoute(func(c *fiber.Ctx) {
		//c.Accepts("application/json")
		id := c.Params("id")
		if id != "" {
			if err := c.JSON(users[id]); err != nil {
				c.SendStatus(500)
			}
			return
		}
		if err := c.JSON(users); err != nil {
			c.SendStatus(500)
		}
	}, "", "/:id")
}

func (u *User) Post() *egowebapi.Route {
	return &egowebapi.Route{
		Path: egowebapi.AddPath(""),
		Handler: func(c *fiber.Ctx) {
			id := c.Query("id")
			u.Lastname = c.Query("lastname")
			u.Firstname = c.Query("firstname")
			users[id] = u
		},
	}
}

func (u *User) Put() *egowebapi.Route {
	return nil
}

func (u *User) Delete() *egowebapi.Route {
	return &egowebapi.Route{
		Path: egowebapi.AddPath("/:id"),
		Handler: func(c *fiber.Ctx) {
			delete(users, c.Params("id"))
		},
	}
}

func (u *User) Options() *egowebapi.Route {
	return &egowebapi.Route{
		Path: egowebapi.AddPath(""),
		Handler: func(ctx *fiber.Ctx) {
			ctx.Append("Allow", "GET, POST, DELETE, OPTIONS")
		},
	}
}
