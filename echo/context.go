package echo

import (
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Context struct {
	Ctx echo.Context
}

func (c *Context) Render(name string, data interface{}, layouts ...string) error {
	return c.Ctx.Render(200, name, data)
}

func (c *Context) Params(key string) string {
	return c.Ctx.Param(key)
}

func (c *Context) Get(key string) string {
	return c.Ctx.Request().Header.Get(key)
}

func (c *Context) Set(key, value string) {
	c.Ctx.Request().Header.Set(key, value)
}

func (c *Context) SendStatus(code int) error {
	return c.Ctx.NoContent(code)
}

func (c *Context) Cookies(key string) string {
	for _, cookie := range c.Ctx.Cookies() {
		if cookie.Name == key {
			return cookie.Value
		}
	}
	return ""
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	c.Ctx.SetCookie(cookie)
}

// TODO ClearCookie
func (c *Context) ClearCookie(key string) {
	for _, cookie := range c.Ctx.Cookies() {
		if cookie.Name == key {

		}
	}
}

func (c *Context) Redirect(location string, status int) error {
	return c.Ctx.Redirect(status, location)
}

func (c *Context) Path() string {
	return c.Ctx.Path()
}

func (c *Context) SendString(code int, s string) error {
	return c.Ctx.String(code, s)
}

func (c *Context) Send(code int, contentType string, b []byte) error {
	return c.Ctx.Blob(code, contentType, b)
}

func (c *Context) SendFile(file string) error {
	return c.Ctx.File(file)
}

func (c *Context) SendStream(code int, contentType string, stream io.Reader) error {
	return c.Ctx.Stream(code, contentType, stream)
}

func (c *Context) JSON(code int, data interface{}) error {
	return c.Ctx.JSON(code, data)
}

func (c *Context) Body() []byte {
	body := c.Ctx.Request().Body
	b, _ := ioutil.ReadAll(body)
	defer body.Close()
	return b
}

func (c *Context) BodyParser(out interface{}) error {
	return c.Ctx.Bind(out)
}

func (c *Context) QueryParam(name string) string {
	return c.Ctx.QueryParam(name)
}

func (c *Context) QueryParams() url.Values {
	return c.Ctx.QueryParams()
}

func (c *Context) Hostname() string {
	c.Ctx.Request()
	return c.Ctx.Request().Host
}

func (c *Context) FormValue(name string) string {
	return c.Ctx.FormValue(name)
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	return c.Ctx.FormFile(name)
}

func (c *Context) Scheme() string {
	return c.Ctx.Scheme()
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.Ctx.MultipartForm()
}
