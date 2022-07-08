package web

import (
	ewa "github.com/ewa-go/ewa"
)

type Home struct{}

func (Home) Get(route *ewa.Route) {
	route.Session()
	route.Handler = func(c *ewa.Context) error {
		return c.Render("home", nil, "layouts/base")
	}
}
