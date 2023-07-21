package web

import (
	"encoding/json"
	"fmt"
	"github.com/ewa-go/ewa"
)

type Login struct{}

func (Login) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return c.Render("login", nil)
	}
}

func (l Login) Post(route *ewa.Route) {
	route.Session(ewa.On).Handler = func(c *ewa.Context) error {

		body := c.Body()
		err := json.Unmarshal(body, &l)
		if err != nil {
			return c.SendString(400, err.Error())
		}

		var sessionId string
		if c.Session != nil {
			sessionId = c.Session.Value
		}

		fmt.Println(sessionId)

		return c.SendString(200, "OK")
	}
}
