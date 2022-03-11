package web

import (
	"errors"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber/src/storage"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Login) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return c.Render("login", nil)
	}
}

func (l Login) Post(route *ewa.Route) {
	route.SetSign(ewa.SignIn)
	route.Handler = func(c *ewa.Context) error {

		err := c.BodyParser(&l)
		if err != nil {
			err = c.SendString(501, err.Error())
			return err
		}

		if l.Username == "user" && l.Password == "Qq123456" {
			storage.SetStorage(c.Identity.SessionId, l.Username)
			return c.SendStatus(200)
		}

		return errors.New("Не верное имя пользователя или пароль!")
	}
}
