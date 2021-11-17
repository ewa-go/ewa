package controllers

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/wsserver"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

type WS struct{}

func (ws *WS) Get(route *ewa.Route) {
	route.SetParams("/:id").WebSocket(ws.Upgrade).Handler = func(c *websocket.Conn) {

		id := c.Params("id")
		wsserver.SetConnection(id, c)

		defer func() {
			err := c.Close()
			if err != nil {
				fmt.Println(err)
			}
			wsserver.DeleteConnection(id)
		}()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}
				return
			}
			log.Printf("messageType: %d, message: %s", mt, msg)

			/*if err := c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}*/
		}
	}
}

func (*WS) Upgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}
