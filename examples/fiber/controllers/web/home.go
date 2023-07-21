package web

import "github.com/ewa-go/ewa"

type Home struct{}

func (Home) Get(route *ewa.Route) {
	route.Session().Handler = func(c *ewa.Context) error {
		return c.Render("home", nil)
	}
}
