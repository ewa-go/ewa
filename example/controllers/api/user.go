package api

import (
	"fmt"
	"github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/src/user"
	"github.com/gofiber/fiber"
)

type User struct {
	*egowebapi.Controller
}

func NewUser(path string) User {

	a := User{
		Controller: egowebapi.NewController("User", "Контроллер для пользователей"),
	}

	path = a.CheckPath(path, a)

	a.SetRoutes(
		egowebapi.NewRoute("GET", path+"/{id}", a.Get, "Вернуть пользователя по id"),
		egowebapi.NewRoute("GET", path, a.Get, "Вернуть пользователей"),
		egowebapi.NewRoute("POST", path, a.Post, "Добавить пользователя"),
		egowebapi.NewRoute("PUT", path, a.Put, "Изменить пользователя"),
		egowebapi.NewRoute("DELETE", path+"/{id}", a.Delete, "Удалить подльзователя"),
	)

	return a
}

func (a User) Get(c *fiber.Ctx) {
	id := c.Params("id")
	if id != "" {
		c.Write(fmt.Sprintf("id: %s, %s", id, user.Get(id).String()))
		return
	}
	s := ""
	for k, v := range user.GetUsers() {
		s += fmt.Sprintf("id: %s, %s\n", k, v.String())
	}
	c.Write(s)
}

func (a User) Post(c *fiber.Ctx) {
	id := c.Query("id")
	lastname := c.Query("lastname")
	firstname := c.Query("firstname")
	user.Set(id, lastname, firstname)
}

func (a User) Put(c *fiber.Ctx) {

}

func (a User) Delete(c *fiber.Ctx) {
	user.Delete(c.Params("id"))
}
