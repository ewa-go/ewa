package web

import (
	"errors"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/web/utils"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *Login) Get(route *ewa.Route) {
	route.SetDescription("Страница Login.html")
	route.Handler = func(ctx *fiber.Ctx) error {
		return ctx.Render("login", nil)
	}
}

func (l *Login) Post(route *ewa.Route) {
	route.SetDescription("Страница Login.html").Login(l.handler, time.Now().Add(24*time.Hour))
}

func (l *Login) handler(ctx *fiber.Ctx, key string) error {

	err := ctx.BodyParser(l)
	if err != nil {
		_, err = ctx.WriteString(err.Error())
		return err
	}

	if l.Username == "user" && l.Password == "Qq123456" {
		storage.SetStorage(key, l.Username)
		utils.SetUser(l.Username)
		return nil
	}

	return errors.New("Не верное имя пользователя или пароль!")
}
