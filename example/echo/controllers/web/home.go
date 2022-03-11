package web

import (
	ewa "github.com/egovorukhin/egowebapi"
)

type Home struct{}

func (Home) Get(route *ewa.Route) {
	route.Session().Permission()
	route.Handler = func(c *ewa.Context) error {
		return c.Render("home", nil, "layouts/base")
	}
}
