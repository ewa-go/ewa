package storage

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber-api/models"
	"github.com/egovorukhin/egowebapi/security"
	"strconv"
)

type User struct{}

func (User) Get(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(true,
		ewa.NewInPath("/{id}", false, "ID пользователя"),
		ewa.NewInQuery("id", false, "ID пользователя"),
		ewa.NewInQueryArray("firstname", "Егор, Вася, Петя", false, "Имена пользователей"),
	)
	route.Handler = func(c *ewa.Context) error {
		id, err := strconv.Atoi(c.Params("id", "0"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		if id > 0 {
			user := models.GetUser(id)
			return c.JSON(200, user)
		}
		users := models.GetUsers()
		return c.JSON(200, users)
	}
	route.SetProduces(ewa.MIMEApplicationJSON)
	route.SetSummary("Get users")
	route.SetDefaultResponse(ewa.NewResponse(ewa.NewSchema(models.User{})).AddHeader("Login", ewa.NewHeader("", false, "Login пользователя")))
	route.SetResponse(200, ewa.NewResponse(ewa.NewSchemaArray(models.User{}), "Return array users"))
	route.SetResponse(422, ewa.NewResponse(nil, "Return parse parameter error"))
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
	route.SetResponse(400, ewa.NewResponse(nil, "Parse body error"))
}

func (User) Put(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(false,
		ewa.NewInQuery("id", false, "id пользователя"),
		ewa.NewInBody(true, ewa.NewSchema(models.User{}), "Необходимо заполнить тело запроса"),
	)
	route.SetProduces(ewa.MIMEApplicationJSON)
	route.SetOperationID("updateUser")
	route.SetSummary("Изменить данные по пользователю")
	route.Handler = func(c *ewa.Context) error {

		id, err := strconv.Atoi(c.QueryParam("id"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		user := models.User{}
		err = c.BodyParser(&user)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		err = user.Update(id)
		if err != nil {
			return c.SendString(400, err.Error())
		}
		return c.SendString(200, "OK")
	}
	route.SetResponse(400, ewa.NewResponse(nil, "Parse body error"))
	route.SetResponse(422, ewa.NewResponse(nil, "Return query error"))
}

func (User) Delete(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(false, ewa.NewInPath("/{id}", true, "ID user"))
	route.SetOperationID("deleteUser")
	route.SetSummary("Удалить пользователя")
	route.Handler = func(c *ewa.Context) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendString(422, err.Error())
		}
		user := models.User{
			Id: id,
		}
		user.Delete()
		return c.SendString(200, "OK")
	}
	route.SetResponse(422, ewa.NewResponse(nil, "Return parse parameter error"))
}
