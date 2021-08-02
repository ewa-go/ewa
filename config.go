package egowebapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/ace"
	"github.com/gofiber/template/amber"
	"github.com/gofiber/template/django"
	"github.com/gofiber/template/handlebars"
	"github.com/gofiber/template/html"
	"github.com/gofiber/template/jet"
	"github.com/gofiber/template/mustache"
	"github.com/gofiber/template/pug"
	"path/filepath"
)

type Config struct {
	Port       int
	Views      *Views
	Secure     *Secure
	Timeout    Timeout
	Static     string
	BasicAuth  *BasicAuth
	Session    *Session
	Permission *Permission
	Store      *session.Config
}

type Secure struct {
	Path string
	Key  string
	Cert string
}

type Views struct {
	Directory string
	Extension Extension
	Layout    string
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
	None       Extension = ""
	Html                 = ".html"
	Ace                  = ".ace"
	Amber                = ".amber"
	Django               = ".django"
	Handlebars           = ".hbs"
	Jet                  = ".jet"
	Mustache             = ".mustache"
	Pug                  = ".pug"
)

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

func (e Extension) Engine(path string, config *Engine) fiber.Views {
	switch e {
	case Html:
		return e.html(path, config)
	case Ace:
		return e.ace(path, config)
	case Amber:
		return e.amber(path, config)
	case Django:
		return e.django(path, config)
	case Handlebars:
		return e.handlebars(path, config)
	case Jet:
		return e.jet(path, config)
	case Mustache:
		return e.mustache(path, config)
	case Pug:
		return e.pug(path, config)
	}
	return nil
}

func (s *Secure) Get() (cert string, key string) {
	key = filepath.Join(s.Path, s.Key)
	cert = filepath.Join(s.Path, s.Cert)
	return cert, key
}

type Timeout struct {
	Read  int
	Write int
	Idle  int
}

func NewTimeout(read, write, idle int) Timeout {
	return Timeout{
		Read:  read,
		Write: write,
		Idle:  idle,
	}
}

func (t Timeout) Get() (read int, write int, idle int) {
	return t.Read, t.Write, t.Idle
}
