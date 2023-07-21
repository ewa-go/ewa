package api

import (
	"github.com/ewa-go/ewa"
	"github.com/ewa-go/ewa/security"
)

type User struct{}

func (User) Get(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth).
		SetParameters(
			ewa.NewPathParam("/{id}", "Идентификатор пользователя"),
		).
		SetEmptyParam("users", "Все пользователи")
	route.Handler = func(c *ewa.Context) error {
		return c.JSON(200, ewa.Map{
			"id":   1,
			"name": "User1",
		})
	}
}

func (User) Post(route *ewa.Route) {
	route.EmptyHandler()
}

func (User) Put(route *ewa.Route) {
	route.EmptyHandler()
}

func (User) Delete(route *ewa.Route) {
	route.EmptyHandler()
}
