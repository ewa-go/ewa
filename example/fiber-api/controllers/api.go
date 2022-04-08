package controllers

import ewa "github.com/egovorukhin/egowebapi"

type Api struct{}

func (Api) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return c.JSON(200, c.Swagger)
	}
}
