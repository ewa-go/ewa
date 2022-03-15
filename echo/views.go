package echo

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"path/filepath"
)

type Renderer struct {
	Root      string
	Extension string
	Layout    string
}

const (
	Html       = ".html"
	Ace        = ".ace"
	Amber      = ".amber"
	Django     = ".django"
	Handlebars = ".hbs"
	Jet        = ".jet"
	Mustache   = ".mustache"
	Pug        = ".pug"
)

// TODO make engine templates
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if name == "" {
		return errors.New("Имя не может быть пустым")
	}

	var files []string
	if r.Layout != "" {
		files = append(files, filepath.Join(r.Root, r.Layout+r.Extension))
	}
	files = append(files, filepath.Join(r.Root, name+r.Extension))

	return template.Must(template.ParseFiles(files...)).ExecuteTemplate(w, name, data)
}

func NewRender(root string, extension string, layout ...string) echo.Renderer {
	r := &Renderer{
		Root:      root,
		Extension: extension,
	}
	if len(layout) > 0 {
		r.Layout = layout[0]
	}
	return r
}
