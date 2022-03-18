package fiber

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Ctx *fiber.Ctx
}

func (c *Context) Render(name string, bind interface{}, layouts ...string) error {
	return c.Ctx.Render(name, bind, layouts...)
}

func (c *Context) Params(key string) string {
	return c.Ctx.Params(key)
}

func (c *Context) Get(key string) string {
	return c.Ctx.Get(key)
}

func (c *Context) Set(key, value string) {
	c.Ctx.Set(key, value)
}

func (c *Context) SendStatus(code int) error {
	return c.Ctx.SendStatus(code)
}

func (c *Context) Cookies(key string) string {
	return c.Ctx.Cookies(key)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	fc := &fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HttpOnly,
		SameSite: strconv.Itoa(int(cookie.SameSite)),
	}
	c.Ctx.Cookie(fc)
}

func (c *Context) ClearCookie(key string) {
	c.Ctx.ClearCookie(key)
}

func (c *Context) Redirect(location string, status int) error {
	return c.Ctx.Redirect(location, status)
}

func (c *Context) Path() string {
	return c.Ctx.Path()
}

func (c *Context) SendString(code int, s string) error {
	return c.Ctx.Status(code).SendString(s)
}

func (c *Context) Send(code int, contentType string, b []byte) error {
	c.Ctx.Set(fiber.HeaderContentType, contentType)
	return c.Ctx.Status(code).Send(b)
}

func (c *Context) SendFile(file string) error {
	return c.Ctx.SendFile(file)
}

func (c *Context) SendStream(code int, contentType string, stream io.Reader) error {
	c.Ctx.Set(fiber.HeaderContentType, contentType)
	return c.Ctx.Status(code).SendStream(stream)
}

func (c *Context) JSON(code int, data interface{}) error {
	return c.Ctx.Status(code).JSON(data)
}

func (c *Context) Body() []byte {
	return c.Ctx.Body()
}

func (c *Context) BodyParser(out interface{}) error {
	return c.Ctx.BodyParser(out)
}

func (c *Context) QueryParam(name string) string {
	return c.Ctx.Query(name)
}

func (c *Context) QueryValues() url.Values {
	values := url.Values{}
	c.Ctx.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		values.Set(string(key), string(value))
	})
	return values
}

func (c *Context) QueryParams(h func(key, value string)) {
	c.Ctx.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		h(string(key), string(value))
	})
}

func (c *Context) Hostname() string {
	return c.Ctx.Hostname()
}

func (c *Context) FormValue(name string) string {
	return c.Ctx.FormValue(name)
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	c.Ctx.Port()
	return c.Ctx.FormFile(name)
}

func (c *Context) Scheme() string {
	return c.Ctx.Protocol()
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.Ctx.MultipartForm()
}
