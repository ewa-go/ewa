package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/models"
	"github.com/egovorukhin/egowebapi/security"
)

type User struct{}

func (User) Get(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(true, ewa.NewInPath("/{id}", false, "ID пользователя"))
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
	route.SetSummary("Вернётся пользовате(ль/ли)")
	route.SetDefaultResponse(ewa.NewResponse(ewa.NewSchema(models.User{})).AddHeader("Login", ewa.NewHeader("", false, "Login пользователя")))
	route.SetOperationID("getUser")
}

func (User) Post(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(false, ewa.NewInBody(true, ewa.NewSchema(models.User{}), "Необходимо заполнить тело запроса"))
	route.SetProduces(ewa.MIMEApplicationJSON)
	route.SetOperationID("setUser")
	route.SetSummary("Добавить пользователя")
	route.Handler = func(c *ewa.Context) error {
		user := models.User{}
		err := c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		user.Set()
		return c.SendString(200, "OK")
	}
}

func (User) Put(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(false, ewa.NewInPath("/{id}", true, "ID пользователя"))
	route.SetProduces(ewa.MIMEApplicationJSON)
	route.SetOperationID("updateUser")
	route.SetSummary("Изменить данные по пользователю")
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
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(false, ewa.NewInPath("/{id}", true, "ID пользователя"))
	route.Handler = func(c *ewa.Context) error {
		user := models.User{
			Id: c.Params("id"),
		}
		user.Delete()
		return c.SendString(200, "OK")
	}
}
