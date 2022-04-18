package egowebapi

import (
	"reflect"
	"regexp"
	"strings"
)

type IGet interface {
	Get(route *Route)
}

type IPost interface {
	Post(route *Route)
}
type IPut interface {
	Put(route *Route)
}

type IDelete interface {
	Delete(route *Route)
}

type IOptions interface {
	Options(route *Route)
}

type IPatch interface {
	Patch(route *Route)
}

type IHead interface {
	Head(route *Route)
}

type ITrace interface {
	Trace(route *Route)
}

type IConnect interface {
	Connect(route *Route)
}

type Controller struct {
	Interface interface{}
	IsShow    bool
	Name      string
	Path      string
	Suffix    []Suffix
	PathTree  []string
	FileTree  []string
	Tag       Tag
}

// SetName Устанавливаем имя контроллера
func (c *Controller) SetName(name string) *Controller {
	c.Name = name
	c.Tag.Name = name
	return c
}

// SetDocs Устанавливаем имя контроллера
func (c *Controller) SetDocs(desc, url string) *Controller {
	c.Tag.ExternalDocs = &ExternalDocs{
		Description: desc,
		URL:         url,
	}
	return c
}

// SetPath Устанавливаем путь контроллера
func (c *Controller) SetPath(path string) *Controller {
	c.Path = path
	return c
}

// SetDescription Устанавливаем описание контроллера
func (c *Controller) SetDescription(desc string) *Controller {
	c.Tag.Description = desc
	return c
}

// SetSuffix Устанавливаем суффикс пути контроллера
func (c *Controller) SetSuffix(suffix ...Suffix) *Controller {
	c.Suffix = append(c.Suffix, suffix...)
	return c
}

// NotShow Установка флага отображения контроллера в swagger
func (c *Controller) NotShow() *Controller {
	c.IsShow = false
	return c
}

// initialize инициализация контролера
func (c *Controller) initialize(basePath string) {

	//Извлекаем имя и путь до "controllers"
	var t reflect.Type
	value := reflect.ValueOf(c.Interface)
	if value.Type().Kind() == reflect.Ptr {
		t = reflect.Indirect(value).Type()
	} else {
		t = value.Type()
	}

	pkg := strings.Replace(
		regexp.MustCompile(`controllers(.*)$`).FindString(t.PkgPath()),
		"controllers",
		"",
		-1,
	)

	if c.Path == "" {
		c.Path = pkg
	}

	c.FileTree = strings.Split(c.Path, "/")
	c.PathTree = c.FileTree
	for _, item := range c.Suffix {
		c.PathTree = insert(c.FileTree, item.Index, item.Value)
	}
	c.Path = strings.Join(c.PathTree, "/")

	if c.Name == "" {
		name := t.Name()
		var path string
		if c.Path != "" && c.Path[:len(basePath)] == basePath {
			path = c.Path[len(basePath):]
		}
		c.Name = strings.ToLower(name)
		c.Tag.Name = strings.ToLower(path + "/" + name)
	}

	c.Path += "/" + c.Name
}

func insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	} else if len(a) < index {
		return a
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
