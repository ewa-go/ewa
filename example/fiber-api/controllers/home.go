package controllers

import ewa "github.com/egovorukhin/egowebapi"

type Home struct {
}

func (Home) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return c.SendString(200, "Home")
	}
}
