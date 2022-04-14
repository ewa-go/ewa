package security

import "fmt"

// Identity Структура описывает идентификацию пользователя
type Identity struct {
	Username string
	AuthName string
}

func (i Identity) String() string {
	return fmt.Sprintf("user: %s, auth_name: %s", i.Username, i.AuthName)
}
