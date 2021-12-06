package web

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
)

type Logout struct{}

func (l *Logout) Get(route *ewa.Route) {
	route.SetDescription("Маршрут /logout")
}

func (l *Logout) Post(route *ewa.Route) {
	route.SetDescription("Маршрут /logout")
	route.Handler = func(ctx *fiber.Ctx, identity *ewa.Identity, key string) error {
		if identity != nil {
			fmt.Println(identity.String())
		}
		storage.DeleteStorage(key)
		return nil
	}
}
