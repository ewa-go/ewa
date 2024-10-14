package security

import (
	"errors"
	"strings"
	"time"
)

type BearerToken struct {
	Param   string
	Handler BearerTokenHandler

	value string
}

type BearerTokenHandler func(token string, isJWT bool) (username string, err error)

func (b *BearerToken) Name() string {
	return BearerTokenAuth
}

func (b *BearerToken) Do() (identity *Identity, err error) {

	if b.value == "" {
		return nil, errors.New("header is required")
	}

	token, jwt, ok := b.parse()
	if !ok {
		return nil, errors.New("invalid token")
	}

	var username string
	if b.Handler != nil {
		username, err = b.Handler(token, jwt)
	}

	identity = &Identity{
		Username: username,
		AuthName: BearerTokenAuth,
		Datetime: time.Now(),
	}

	return
}

func (b *BearerToken) Definition() Definition {
	return Definition{
		Type:        TypeBearerToken,
		Description: "Bearer Token Authorization",
	}
}

func (b *BearerToken) SetValue(value string) *BearerToken {
	b.value = value
	return b
}

func (b *BearerToken) parse() (token string, isJWT, ok bool) {
	const prefix = "Bearer "
	if b.value[:len(prefix)] == prefix {
		token = b.value[len(prefix):]
		ok = true
		i := strings.Index(b.value, ".")
		if i > 0 {
			isJWT = true
		}
	}
	return
}
