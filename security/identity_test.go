package security

import (
	"fmt"
	"testing"
	"time"
)

func newIdentity() *Identity {
	return &Identity{
		Username: "username",
		AuthName: "Session",
		Datetime: time.Now(),
	}
}

func TestIdentity_SetVariable(t *testing.T) {
	i := newIdentity()
	fmt.Println(i.String())
	i.SetVariable("is_admin", true)
	fmt.Println(i.String())
}

func TestIdentity_SetVariables(t *testing.T) {
	i := newIdentity()
	fmt.Println(i.String())
	i.SetVariables(map[string]interface{}{"is_admin": false})
	fmt.Println(i.String())
}
