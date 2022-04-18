package egowebapi

import (
	"fmt"
	"testing"
)

type User struct{}

func (User) Get(route *Route) {

}

func TestInitialize(t *testing.T) {

	c := Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/api/storage",
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)

	hostname := Suffix{
		Index: 3,
		Value: "{hostname}",
	}

	c = Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/api/storage",
		Suffix: []Suffix{
			hostname,
		},
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)
}
