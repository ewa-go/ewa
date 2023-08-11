package security

import (
	"fmt"
	"time"
)

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	Username string    `json:"username"`
	AuthName string    `json:"auth_name"`
	Datetime time.Time `json:"datetime"`
}

func (i Identity) String() string {
	return fmt.Sprintf("user: %s, auth_name: %s", i.Username, i.AuthName)
}
