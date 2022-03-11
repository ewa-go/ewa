package egowebapi

import (
	"net/http"
	"net/url"
)

type Context struct {
	Identity *Identity
	IContext
}

type IContext interface {
	Render(name string, data interface{}, layouts ...string) error
	Params(key string) string
	Get(key string) string
	Set(key string, value string)
	SendStatus(code int) error
	Send(code int, contentType string, b []byte) error
	SendString(code int, s string) error
	Cookies(key string) string
	SetCookie(cookie *http.Cookie)
	ClearCookie(key string)
	Redirect(location string, status int) error
	Path() string
	JSON(code int, data interface{}) error
	Body() []byte
	BodyParser(out interface{}) error
	QueryParam(name string) string
	QueryParams() url.Values
}

func NewContext(c IContext) *Context {
	return &Context{IContext: c}
}
