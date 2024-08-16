package ewa

import (
	"fmt"
	"testing"
)

type User struct{}

func (User) Get(route *Route) {
	route.Handler = func(c *Context) error {
		return nil
	}
}

func TestInitialize(t *testing.T) {

	c := Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/api/storage",
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)

	hostname := &Suffix{
		Index: 3,
		Value: "{hostname}",
	}
	id := &Suffix{
		Index: 4,
		Value: "{id}",
	}

	c = Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/api/storage",
		Suffix: []*Suffix{
			hostname,
			id,
		},
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)

	c = Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/",
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)

	c = Controller{
		Interface: new(User),
		IsShow:    true,
		Path:      "/",
		Name:      "MyUser",
	}

	c.initialize("/api")
	fmt.Println(c.Name, c.Path, c.Tag.Name)
}
