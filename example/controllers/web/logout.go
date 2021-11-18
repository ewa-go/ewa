package web

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
)

type Logout struct{}

func (l *Logout) Get(route *ewa.Route) {
	route.SetDescription("Маршрут /logout")
}

func (l *Logout) Post(route *ewa.Route) {
	route.SetDescription("Маршрут /logout").
		WebAuth(false).Handler = func(ctx *fiber.Ctx, key string) error {
		storage.DeleteStorage(key)
		return nil
	}
}
