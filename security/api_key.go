package security

import (
	"errors"
	"fmt"
)

type ApiKey struct {
	KeyName string
	Param   string
	value   string
	Handler ApiKeyAuthHandler
}

type ApiKeyAuthHandler func(token string) (username string, err error)

const (
	ParamQuery  = "path"
	ParamHeader = "header"
)

func (a *ApiKey) SetValue(value string) *ApiKey {
	a.value = value
	return a
}

func (a ApiKey) Do() (identity *Identity, err error) {

	/*var value string
	switch a.Param {
	// Пытаемся получить из заголовка токен
	case ParamQuery:
		value = c.QueryParam(a.KeyName)
		break
	// Если не нашли в заголовке, то ищем в переменных запроса адресной строки
	case ParamHeader:
		value = c.Get(a.KeyName)
		break
	}*/

	if a.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", a.Param))
	}

	username := ""
	if a.Handler != nil {
		username, err = a.Handler(a.value)
	}

	identity = &Identity{
		Username: username,
		AuthName: ApiKeyAuth,
	}

	return
}

func (a ApiKey) Definition() Definition {
	return Definition{
		Type:        TypeApiKey,
		Description: fmt.Sprintf("Api Key Authorization. Set name: %s, parameter: %s", a.KeyName, a.Param),
	}
}
