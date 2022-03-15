package websocket

import ewa "github.com/egovorukhin/egowebapi"

type Conn struct {
	IConn
}

type IConn interface {
	ReadMessage() (int, []byte, error)
	ReadJSON(data interface{}) error
	WriteMessage(messageType int, data []byte) error
	WriteJSON(v interface{}) error
	Close() error
}

func New(h func(conn *Conn)) ewa.Handler {
	return func(c *ewa.Context) error {
		h(c)
	}
}
