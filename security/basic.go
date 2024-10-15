package security

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

type Basic struct {
	header  string
	Handler BasicAuthHandler
}

type BasicAuthHandler func(user string, pass string) error

func (b *Basic) parse() (username, password string, ok bool) {
	const prefix = "Basic "
	if len(b.header) < len(prefix) || !strings.EqualFold(b.header[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(b.header[len(prefix):])
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

func (b *Basic) SetHeader(header string) *Basic {
	b.header = header
	return b
}

func (b *Basic) Name() string {
	return BasicAuth
}

func (b *Basic) Do() (*Identity, error) {

	err := errors.New(`basic realm="Необходимо указать имя пользователя и пароль"`)
	if b.header == "" {
		return nil, err
	}

	username, password, ok := b.parse()
	if !ok {
		return nil, err
	}

	err = b.Handler(username, password)
	if err != nil {
		return nil, err
	}

	identity := &Identity{
		Username: username,
		AuthName: BasicAuth,
		Datetime: time.Now(),
	}

	return identity, nil
}

func (b *Basic) Definition() Definition {
	return Definition{
		Type:        TypeBasic,
		Description: "Basic Authorization",
	}
}
