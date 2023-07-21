package web

import "github.com/ewa-go/ewa"

type Logout struct{}

func (Logout) Get(route *ewa.Route) {
	route.Session(ewa.Off).Handler = func(c *ewa.Context) error {
		return nil
	}
}
