package web

import (
	"fmt"
	ewa "github.com/ewa-go/ewa"
	"github.com/ewa-go/ewa/example/echo/src/storage"
)

type Logout struct{}

func (Logout) Get(route *ewa.Route) {
	route.SetSign(ewa.SignOut)
	route.Handler = func(c *ewa.Context) error {
		if c.Identity != nil {
			fmt.Println(c.Identity.String())
			storage.DeleteStorage(c.Identity.SessionId)
		}
		return c.SendStatus(200)
	}
}
