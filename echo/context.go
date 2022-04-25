package echo

import (
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Context struct {
	Ctx echo.Context
}

func (c *Context) Render(name string, data interface{}, layouts ...string) error {
	return c.Ctx.Render(200, name, data)
}

func (c *Context) Params(key string, defaultValue ...string) string {
	value := c.Ctx.Param(key)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
}

func (c *Context) Get(key string, defaultValue ...string) string {
	value := c.Ctx.Request().Header.Get(key)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
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

func (c *Context) SaveFile(fh *multipart.FileHeader, path string) (err error) {
	var (
		f  multipart.File
		ff *os.File
	)
	f, err = fh.Open()
	if err != nil {
		return
	}

	var ok bool
	if ff, ok = f.(*os.File); ok {
		// Windows can't rename files that are opened.
		if err = f.Close(); err != nil {
			return
		}

		// If renaming fails we try the normal copying method.
		// Renaming could fail if the files are on different devices.
		if os.Rename(ff.Name(), path) == nil {
			return nil
		}

		// Reopen f for the code below.
		if f, err = fh.Open(); err != nil {
			return
		}
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	if ff, err = os.Create(path); err != nil {
		return
	}
	defer func() {
		e := ff.Close()
		if err == nil {
			err = e
		}
	}()

	return
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

func (c *Context) QueryParam(name string, defaultValue ...string) string {
	value := c.Ctx.QueryParam(name)
	if value == "" && defaultValue != nil {
		return defaultValue[0]
	}
	return value
}

func (c *Context) QueryValues() url.Values {
	return c.Ctx.QueryParams()
}

func (c *Context) QueryParams(h func(key, value string)) {
	for k, v := range c.Ctx.QueryParams() {
		s := ""
		if len(v) > 0 {
			s = v[0]
		}
		h(k, s)
	}
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
