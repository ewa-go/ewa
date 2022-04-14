package security

import (
	"encoding/base64"
	"errors"
	"strings"
)

type Basic struct {
	header  string
	Handler BasicAuthHandler
}

type BasicAuthHandler func(user string, pass string) bool

func (b Basic) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	i := strings.IndexByte(cs, ':')
	if i < 0 {
		return
	}
	return cs[:i], cs[i+1:], true
}

func (b *Basic) SetHeader(header string) {
	b.header = header
}

func (b Basic) Do() (*Identity, error) {

	err := errors.New(`Basic realm="Необходимо указать имя пользователя и пароль"`)
	if b.header == "" {
		return nil, err
	}

	username, password, ok := b.parseBasicAuth(b.header)
	if !ok || !b.Handler(username, password) {
		return nil, err
	}

	identity := &Identity{
		Username: username,
		AuthName: BasicAuth,
	}

	return identity, nil
}

func (b Basic) Definition() Definition {
	return Definition{
		Type:        TypeBasic,
		Description: "Basic Authorization",
	}
}
