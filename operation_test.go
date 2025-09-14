package ewa

import (
	"fmt"
	"testing"
)

func TestGetPathParams(t *testing.T) {
	o := Operation{}
	o.Parameters = append(o.Parameters, NewPathParam("/{id}", "ID"))
	o.Parameters = append(o.Parameters, NewPathParam("/{id}/config", "Session"))
	o.Parameters = append(o.Parameters, NewPathParam("/{id}/config/{config_id}", "ConfigID"))
	params := o.getPathParams()
	fmt.Println(params)
}

func TestGetParams(t *testing.T) {
	o := Operation{}
	o.Parameters = append(o.Parameters, NewPathParam("/{id}", "ID"))
	o.Parameters = append(o.Parameters, NewPathParam("/{id}/config", "Session"))
	o.Parameters = append(o.Parameters, NewPathParam("/{id}/config/{config_id}", "ConfigID"))
	params := o.getParams()
	fmt.Println(params)
}
