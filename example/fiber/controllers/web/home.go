package web

import (
	ewa "github.com/egovorukhin/egowebapi"
)

type Home struct{}

func (Home) Get(route *ewa.Route) {
	route.Auth(ewa.SessionAuth).Permission()
	route.Handler = func(c *ewa.Context) error {
		return c.ViewRender(nil)
	}
}
