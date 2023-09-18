package security

import (
	"fmt"
	"time"
)

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	Username  string                 `json:"username"`
	AuthName  string                 `json:"auth_name"`
	Datetime  time.Time              `json:"datetime"`
	Variables map[string]interface{} `json:"variables"`
}

func (i *Identity) GetVariable(name string) interface{} {
	if v, ok := i.Variables[name]; ok {
		return v
	}
	return nil
}

func (i *Identity) SetVariable(name string, value interface{}) *Identity {
	if i.Variables == nil {
		i.Variables = make(map[string]interface{})
	}
	i.Variables[name] = value
	return i
}

func (i *Identity) DeleteVariable(name string) *Identity {
	if _, ok := i.Variables[name]; ok {
		delete(i.Variables, name)
	}
	return i
}

func (i *Identity) String() string {
	return fmt.Sprintf("user: %s, auth_name: %s, datetime: %s, variables: %v", i.Username, i.AuthName, i.Datetime, i.Variables)
}
