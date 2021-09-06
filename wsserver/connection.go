package wsserver

import (
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/websocket/v2"
	"sync"
)

type Connection struct {
	conn *websocket.Conn
}

type Connections struct {
	sync.Map
}

func (c *Connections) Add(id interface{}, conn *Connection) string {
	uuid := utils.UUID()
	c.Store(uuid, conn)
	return uuid
}
