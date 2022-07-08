package web

import ewa "github.com/ewa-go/ewa"

type Index struct{}

func (Index) Get(route *ewa.Route) {
	//route.Session()
	route.Handler = func(c *ewa.Context) error {
		return c.Render("index", nil)
	}
}
