package security

import (
	"errors"
	"fmt"
	"time"
)

type ApiKey struct {
	KeyName string
	Param   string
	value   string
	Handler ApiKeyAuthHandler
}

type ApiKeyAuthHandler func(token string) (username string, err error)

func (a *ApiKey) SetValue(value string) *ApiKey {
	a.value = value
	return a
}

func (a *ApiKey) Name() string {
	return ApiKeyAuth
}

func (a *ApiKey) Do() (identity *Identity, err error) {

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
		Datetime: time.Now(),
	}

	return
}

func (a *ApiKey) Definition() Definition {
	return Definition{
		Type:        TypeApiKey,
		In:          a.Param,
		Name:        a.KeyName,
		Description: fmt.Sprintf("Api Key Authorization. Set name: %s, in: %s", a.KeyName, a.Param),
	}
}
