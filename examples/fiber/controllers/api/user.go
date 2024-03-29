package api

import (
	"fmt"
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
		req := c.HttpRequest()
		fmt.Println(req)
		var isAdmin bool
		if c.Identity != nil {
			if value, ok := c.Identity.Variables["is_admin"]; ok {
				isAdmin = value.(bool)
			}
		}
		return c.JSON(200, ewa.Map{
			"id":       1,
			"name":     "User1",
			"is_admin": isAdmin,
		})
	}
}

func (User) Post(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth).Handler = func(c *ewa.Context) error {
		return c.JSON(200, ewa.Map{
			"ok": true,
		})
	}
}

func (User) Put(route *ewa.Route) {
	route.EmptyHandler()
}

func (User) Delete(route *ewa.Route) {
	route.EmptyHandler()
}
