package wsserver

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type Connections struct {
	sync.Map
}

var conns *Connections

func SetConnection(id string, conn *websocket.Conn) string {
	if conns == nil {
		conns = &Connections{}
	}
	conns.Store(id, conn)
	return id
}

func GetConnection(id string) *websocket.Conn {
	if conns == nil {
		return nil
	}
	conn, _ := conns.Load(id)
	return conn.(*websocket.Conn)
}

func DeleteConnection(id string) {
	if conns == nil {
		return
	}
	conns.Delete(id)
}
