package web

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber/src/storage"
)

type Logout struct{}

func (Logout) Get(route *ewa.Route) {
	route.SetSign(ewa.SignOut)
	route.Handler = func(c *ewa.Context) error {
		if c.SessionId != nil {
			sessionId := c.SessionId.(string)
			fmt.Println(sessionId)
			storage.DeleteStorage(sessionId)
		}
		return c.SendStatus(200)
	}
}
