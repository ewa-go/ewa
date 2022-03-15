package controllers

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber/src/wsserver"
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

type WS struct{}

func (ws *WS) Get(route *ewa.Route) {
	route.SetParams("/:id")
	route.Handler = websocket.New() func(c *ewa.Context) error {

		conn := c.WebSocket

		id := c.Params("id")
		wsserver.AddClient(&wsserver.Client{
			Id:      id,
			Conn:    conn,
			Created: time.Now(),
		})

		defer func() {
			err := conn.Close()
			if err != nil {
				fmt.Println(err)
			}
			wsserver.DeleteClient(id)
		}()

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			log.Printf("messageType: %d, message: %s", mt, msg)

			if err := conn.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

		return nil
	}
	/*route.Handler = func(c *websocket.Conn) {

		id := c.Params("id")
		wsserver.AddClient(&wsserver.Client{
			Id:      id,
			Conn:    c,
			Created: time.Now(),
		})

		defer func() {
			err := c.Close()
			if err != nil {
				fmt.Println(err)
			}
			wsserver.DeleteClient(id)
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

			if err := c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}*/
}

/*func (WS) Upgrade(c *ewa.Context) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}*/
