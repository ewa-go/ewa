package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

var users = map[string]*User{}

type User struct {
	Lastname  string
	Firstname string
}

func (u *User) Get() *ewa.Route {
	return ewa.NewRoute(

		func(c *fiber.Ctx) error {

			//c.Accepts("application/json")
			id := c.Params("id")
			if id != "" {
				if err := c.JSON(users[id]); err != nil {
					c.SendStatus(500)
					return err
				}
				return nil
			}
			if err := c.JSON(users); err != nil {
				c.SendStatus(500)
				return err
			}

			return nil

		},
		"", "/:id")
}

func (u *User) Post() *ewa.Route {
	return &ewa.Route{
		Path: ewa.AddPath(""),
		Handler: func(c *fiber.Ctx) error {
			ba := c.Get("Authorization")
			id := c.Query("id")
			u.Lastname = c.Query("lastname")
			u.Firstname = c.Query("firstname")
			users[id] = u
			return nil
		},
	}
}

func (u *User) Put() *ewa.Route {
	return nil
}

func (u *User) Delete() *ewa.Route {
	return &ewa.Route{
		Path: ewa.AddPath("/:id"),
		Handler: func(c *fiber.Ctx) error {
			delete(users, c.Params("id"))
			return nil
		},
	}
}

func (u *User) Options() *ewa.Route {
	return &ewa.Route{
		Path: ewa.AddPath(""),
		Handler: func(ctx *fiber.Ctx) error {
			ctx.Append("Allow", "GET, POST, DELETE, OPTIONS")
			return nil
		},
	}
}
