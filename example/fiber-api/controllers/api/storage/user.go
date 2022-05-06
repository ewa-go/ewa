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
	route.SetEmptyParam(
		ewa.NewEmptyPathParam("Get users").
			SetResponse(200, ewa.NewResponse(ewa.NewSchemaArray(models.User{}), "Return array users")),
	)
	route.SetParameters(
		ewa.NewPathParam("/{id}", "ID users").SetType(ewa.TypeInteger),
		ewa.NewQueryParam("id", false, "ID users"),
		ewa.NewQueryArrayParam("firstname", "User1, User2, User3", false, "Name users"),
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
	route.SetSummary("Get user")
	route.SetDefaultResponse(ewa.NewResponse(ewa.NewSchema(models.User{})).AddHeader("Login", ewa.NewHeader("", false, "User login")))
	route.SetResponse(422, ewa.NewResponse(nil, "Return parse parameter error"))
	route.SetResponse(200, ewa.NewResponse(ewa.NewSchema(models.User{}), "Return user struct"))
}

func (User) Post(route *ewa.Route) {
	route.SetSecurity(security.BasicAuth)
	route.SetParameters(ewa.NewBodyParam(true, ewa.NewSchema(models.User{}), "Must have request body"))
	route.SetSummary("Create user")
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
	route.SetParameters(
		ewa.NewQueryParam("id", false, "id user"),
		ewa.NewBodyParam(true, ewa.NewSchema(models.User{}), "Must have request body"),
	)
	route.SetSummary("Update user")
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
	route.SetParameters(ewa.NewPathParam("/{id}", "ID user"))
	route.SetSummary("Delete user")
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
