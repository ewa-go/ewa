package controllers

/*type WS struct{}

func (ws *WS) Get(route *ewa.Route) {
	route.SetParams("/:id")
	route.Handler = websocket.New(func(conn *websocket.Conn) {
		id := conn.Params("id")
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
				log.Println(err)
				break
			}
			log.Printf("messageType: %d, message: %s", mt, msg)

			if err := conn.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	})
}*/

/*func (WS) Upgrade(c *ewa.Context) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}*/
