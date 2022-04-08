package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/models"
	"github.com/egovorukhin/egowebapi/swagger/v2"
)

type User struct{}

func (User) Get(route *ewa.Route) {
	route.SetParams("", "/:id").Auth(ewa.BasicAuth)
	route.Handler = func(c *ewa.Context) error {
		id := c.Params("id")
		if id != "" {
			user := models.GetUser(id)
			return c.JSON(200, user)
		}
		users := models.GetUsers()
		return c.JSON(200, users)
	}
	route.SetProduces(ewa.MIMEApplicationJSON)
	route.SetParameters(v2.Parameter{
		Name:     "id",
		In:       v2.ParameterTypePath,
		Required: false,
		Type:     "integer",
	})
	route.SetResponse(200, v2.Response{
		Description: "Возвращает - OK",
	})
}

func (User) Post(route *ewa.Route) {
	route.Auth(ewa.BasicAuth)
	route.Handler = func(c *ewa.Context) error {
		user := models.User{}
		err := c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		user.Set()
		return c.SendString(200, "OK")
	}
	route.SetParameters(v2.Parameter{
		Name:     v2.ParameterTypeBody,
		In:       v2.ParameterTypeBody,
		Required: true,
		Body:     models.User{},
	})
}

func (User) Put(route *ewa.Route) {
	route.SetParams("/:id").Auth(ewa.BasicAuth)
	route.Handler = func(c *ewa.Context) error {
		user := models.User{}
		err := c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		err = user.Update(c.QueryParam("id"))
		if err != nil {
			return c.SendString(400, err.Error())
		}
		return c.SendString(200, "OK")
	}
}

func (User) Delete(route *ewa.Route) {
	route.Auth(ewa.BasicAuth)
	route.Handler = func(c *ewa.Context) error {
		user := models.User{
			Id: c.Params("id"),
		}
		user.Delete()
		return c.SendString(200, "OK")
	}
}

func (User) Tag() v2.Tag {
	return v2.Tag{
		Description: "Описание",
	}
}
