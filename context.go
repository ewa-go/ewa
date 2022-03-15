package egowebapi

import (
	"github.com/egovorukhin/egowebapi/websocket"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Context struct {
	Identity  *Identity
	View      *View
	WebSocket *websocket.Conn
	IContext
}

type View struct {
	Filename string
	Filepath string
	Layout   string
}

type IContext interface {
	Render(name string, data interface{}, layouts ...string) error
	Params(key string) string
	Get(key string) string
	Set(key string, value string)
	SendStatus(code int) error
	Send(code int, contentType string, b []byte) error
	SendString(code int, s string) error
	SendFile(file string) error
	SendStream(code int, contentType string, stream io.Reader) error
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
	Hostname() string
	FormValue(name string) string
	FormFile(name string) (*multipart.FileHeader, error)
	Scheme() string
	MultipartForm() (*multipart.Form, error)
}

func NewContext(c IContext) *Context {
	return &Context{IContext: c}
}

func (c *Context) ViewRender(data interface{}) error {
	return c.Render(c.View.Filename, data, c.View.Layout)
}
