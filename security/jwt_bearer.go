package security

import (
	"errors"
	"fmt"
	"time"
)

type JWTBearer struct {
	KeyName string
	Param   string
	value   string
	Handler JWTBearerAuthHandler
}

type JWTBearerAuthHandler func(token string) (username string, err error)

func (*JWTBearer) SetValues(v Values) {

}

func (a *JWTBearer) Do() (identity *Identity, err error) {

	if a.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", a.Param))
	}

	username := ""
	if a.Handler != nil {
		username, err = a.Handler(a.value)
	}

	identity = &Identity{
		Username: username,
		AuthName: JWTBearerAuth,
		Datetime: time.Now(),
	}

	return
}

func (a *JWTBearer) Definition() Definition {
	return Definition{
		Type:        TypeJWTBearer,
		In:          a.Param,
		Name:        a.KeyName,
		Description: fmt.Sprintf("JWTBearer Authorization. Set name: %s, in: %s", a.KeyName, a.Param),
	}
}
