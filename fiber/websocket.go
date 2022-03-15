package fiber

import "github.com/gofiber/websocket/v2"

type Conn struct {
	C *websocket.Conn
}

func (c *Conn) ReadMessage() (int, []byte, error) {
	return c.C.ReadMessage()
}

func (c *Conn) ReadJSON(data interface{}) error {
	return c.C.ReadJSON(data)
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	return c.C.WriteMessage(messageType, data)
}

func (c *Conn) WriteJSON(v interface{}) error {
	return c.C.WriteJSON(v)
}

func (c *Conn) Close() error {
	return c.C.Close()
}
