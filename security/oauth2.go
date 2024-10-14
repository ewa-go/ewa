package security

import (
	"errors"
	"fmt"
	"time"
)

type OAuth2 struct {
	HeaderPrefix string
	Param        string
	Flow         Flow
	Handler      OAuth2AuthHandler

	value string
}

type Flow struct {
	Type   string
	Url    string
	Scopes map[string]string
}

const (
	FlowImplicit    = "implicit"
	FlowPassword    = "password"
	FlowAccessCode  = "access_code"
	FlowApplication = "application"
)

type OAuth2AuthHandler func(token string) (username string, err error)

func (o *OAuth2) SetValue(value string) {
	o.value = value
}

func (o *OAuth2) Name() string {
	return OAuth2Auth
}

func (o *OAuth2) Do() (identity *Identity, err error) {

	if o.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", o.Param))
	}

	username := ""
	if o.Handler != nil {
		username, err = o.Handler(o.value)
	}

	identity = &Identity{
		Username: username,
		AuthName: OAuth2Auth,
		Datetime: time.Now(),
	}

	return
}

func (o *OAuth2) Definition() (d Definition) {
	d.Type = TypeOAuth2
	d.Description = "OAuth2 Authorization"
	d.Flow = o.Flow.Type
	d.Scopes = o.Flow.Scopes
	switch o.Flow.Type {
	case FlowImplicit, FlowAccessCode:
		d.AuthorizationURL = o.Flow.Url
	case FlowPassword, FlowApplication:
		d.TokenURL = o.Flow.Url
	}
	return
}

func (o *OAuth2) parse() (token string, ok bool) {
	switch o.Param {
	case ParamHeader:
		if o.value[:len(o.HeaderPrefix)] == o.HeaderPrefix {
			return o.value[len(o.HeaderPrefix):], true
		}
	case ParamQuery:
		return o.value, true
	}
	return
}
