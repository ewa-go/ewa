package wsserver

import (
	"github.com/egovorukhin/egowebapi/websocket"
	"sync"
	"time"
)

type Client struct {
	Id      string          `json:"id"`
	Conn    *websocket.Conn `json:"-"`
	Created time.Time       `json:"created"`
}

type Clients struct {
	sync.Map
}

var clients *Clients

func AddClient(client *Client) {
	if clients == nil {
		clients = new(Clients)
	}
	clients.Store(client.Id, client)
}

func GetClient(id string) *Client {
	if clients == nil {
		return nil
	}
	value, _ := clients.Load(id)
	return value.(*Client)
}

func GetClients() (c []*Client) {
	if clients == nil {
		return nil
	}
	clients.Range(func(key, value interface{}) bool {
		c = append(c, value.(*Client))
		return true
	})
	return
}

func DeleteClient(id string) {
	if clients == nil {
		return
	}
	clients.Delete(id)
}
