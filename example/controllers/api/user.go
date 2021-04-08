package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

var users = Users{}

type User struct {
	Id        string
	Lastname  string
	Firstname string
}

type Users []User

func (u *User) Get() *ewa.Route {
	return ewa.NewRoute(

		func(c *fiber.Ctx) error {

			id := c.Params("id")
			if id != "" {
				_, user := GetUser(id)
				if err := c.JSON(user); err != nil {
					c.SendStatus(500)
					return err
				}
				return nil
			}
			if err := c.JSON(GetUsers()); err != nil {
				c.SendStatus(500)
				return err
			}

			return nil

		},
		"", "/:id")
}

func (u *User) Post() *ewa.Route {
	return &ewa.Route{
		Path: nil,
		Handler: func(c *fiber.Ctx) error {
			user := &User{}
			user.Id = c.Query("id")
			user.Lastname = c.Query("lastname")
			user.Firstname = c.Query("firstname")
			SetUser(*user)
			return nil
		},
	}
}

func (u *User) Put() *ewa.Route {
	return &ewa.Route{
		Path: ewa.AddPath("/:id"),
		Handler: func(c *fiber.Ctx) error {
			u.Id = c.Params("id")
			u.Update()
			return nil
		},
	}
}

func (u *User) Delete() *ewa.Route {
	return &ewa.Route{
		Path: ewa.AddPath("/:id"),
		Handler: func(c *fiber.Ctx) error {
			u.Id = c.Params("id")
			u.Remove()
			return nil
		},
	}
}

func (u *User) Options() *ewa.Route {
	return &ewa.Route{
		Path: nil,
		Handler: func(ctx *fiber.Ctx) error {
			ctx.Append("Allow", "GET, POST, DELETE, OPTIONS")
			return nil
		},
	}
}

func GetUsers() Users {
	return users
}

func GetUser(id string) (int, *User) {
	for i, user := range users {
		if user.Id == id {
			return i, &user
		}
	}
	return -1, nil
}

func SetUser(u User) {
	users = append(users, u)
}

func (u *User) Update() {
	i, _ := GetUser(u.Id)
	users[i] = *u
}

func (u *User) Remove() {
	i, _ := GetUser(u.Id)
	users = append(users[:i], users[i+1:]...)
}
