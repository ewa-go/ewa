package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/ace"
	"github.com/gofiber/template/amber"
	"github.com/gofiber/template/django"
	"github.com/gofiber/template/handlebars"
	"github.com/gofiber/template/html"
	"github.com/gofiber/template/jet"
	"github.com/gofiber/template/mustache"
	"github.com/gofiber/template/pug"
)

type Views struct {
	Directory string
	Extension Extension
	Engine    *Engine
}

type Engine struct {
	Reload bool
	Debug  bool
	Layout string
	Delims *Delims
}

type Delims struct {
	Left  string
	Right string
}

type Extension string

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

const (
	defaultBufferSize = 1024 * 4
)

func NewViews(dir string, ext Extension, e *Engine) fiber.Views {
	return Views{
		Directory: dir,
		Extension: ext,
		Engine:    e,
	}.engine()
}

func (v *Views) SetEngine(e *Engine) *Views {
	v.Engine = e
	return v
}

// html template
func (e Extension) html(path string, config *Engine) *html.Engine {
	engine := html.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// ace template
func (e Extension) ace(path string, config *Engine) *ace.Engine {
	engine := ace.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// amber template
func (e Extension) amber(path string, config *Engine) *amber.Engine {
	engine := amber.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// django template
func (e Extension) django(path string, config *Engine) *django.Engine {
	engine := django.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// handlebars template
func (e Extension) handlebars(path string, config *Engine) *handlebars.Engine {
	engine := handlebars.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// jet template
func (e Extension) jet(path string, config *Engine) *jet.Engine {
	engine := jet.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// mustache template
func (e Extension) mustache(path string, config *Engine) *mustache.Engine {
	engine := mustache.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

// pug template
func (e Extension) pug(path string, config *Engine) *pug.Engine {
	engine := pug.New(path, string(e))
	if config != nil {
		engine.Reload(config.Reload)
		engine.Debug(config.Debug)
		if config.Delims != nil {
			engine.Delims(config.Delims.Left, config.Delims.Right)
		}
		if config.Layout != "" {
			engine.Layout(config.Layout)
		}
	}
	return engine
}

func (v Views) engine() fiber.Views {
	e := v.Extension
	switch e {
	case Html:
		return e.html(v.Directory, v.Engine)
	case Ace:
		return e.ace(v.Directory, v.Engine)
	case Amber:
		return e.amber(v.Directory, v.Engine)
	case Django:
		return e.django(v.Directory, v.Engine)
	case Handlebars:
		return e.handlebars(v.Directory, v.Engine)
	case Jet:
		return e.jet(v.Directory, v.Engine)
	case Mustache:
		return e.mustache(v.Directory, v.Engine)
	case Pug:
		return e.pug(v.Directory, v.Engine)
	}
	return nil
}
