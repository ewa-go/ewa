package api

import (
	ewa "github.com/ewa-go/ewa"
)

type Messages struct {
}

func (u *Messages) Get(route *ewa.Route) {
	route.EmptyHandler()
}
