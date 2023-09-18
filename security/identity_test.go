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

func TestIdentity_GetVariable(t *testing.T) {
	i := newIdentity()
	i.Variables = make(map[string]interface{})
	i.Variables["is_admin"] = true
	fmt.Println(i.GetVariable("is_admin"))
}

func TestIdentity_DeleteVariable(t *testing.T) {
	i := newIdentity()
	i.Variables = make(map[string]interface{})
	i.Variables["is_admin"] = true
	fmt.Println(i.String())
	i.DeleteVariable("is_admin")
	fmt.Println(i.String())
}
