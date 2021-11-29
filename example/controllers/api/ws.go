package api

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/src/wsserver"
	"github.com/gofiber/fiber/v2"
)

type WS struct {
	Id      string
	Message string `json:"message"`
}

func (*WS) Get(route *ewa.Route) {
	route.SetParams("", "/:id")
	route.Handler = func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		if len(id) > 0 {
			return ctx.Status(200).JSON(wsserver.GetClient(id))
		}
		return ctx.Status(200).JSON(wsserver.GetClients())
	}
}

func (ws *WS) Post(route *ewa.Route) {
	route.SetParams("", "/:id")
	route.Handler = func(ctx *fiber.Ctx) error {
		err := ctx.BodyParser(&ws)
		if err != nil {
			return ctx.Status(501).SendString(err.Error())
		}
		id := ctx.Params("id")
		if len(id) == 0 {
			id = ws.Id
		}
		client := wsserver.GetClient(id)
		err = client.Conn.WriteJSON(ws.Message)
		if err != nil {
			return ctx.Status(501).SendString(err.Error())
		}
		return ctx.Status(200).SendString("Done")
	}
}
