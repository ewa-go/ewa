package ewa

import (
	"fmt"
	"testing"
)

func TestGetPathParams(t *testing.T) {
	o := Operation{}
	params := o.getPathParams()
	fmt.Println(params)
}
