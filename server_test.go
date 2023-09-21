package ewa

import (
	"fmt"
	"testing"
)

type Test struct{}

func (Test) Get(route *Route) {
	route.Handler = func(c *Context) error {
		return nil
	}
}

func TestNewSuffix(t *testing.T) {

	hostname := Suffix{
		Index:       2,
		Value:       "hostname",
		isParam:     false,
		Description: "Hostname",
	}

	client := Suffix{
		Index:       2,
		Value:       "client",
		isParam:     false,
		Description: "Client",
	}

	s := NewSuffix(hostname, client)
	fmt.Printf("%#v", s)
}

func TestAdd(t *testing.T) {
	s := &Server{}
	s.Register(new(Test))
	err := s.Start()
	if err != nil {
		t.Fatal(err)
	}
}
