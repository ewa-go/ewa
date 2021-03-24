package egowebapi

import (
	"reflect"
	"regexp"
	"strings"
)

type Controller struct {
	Name        string
	Description string
	Routes      Routes
}

type Controllers []*Controller

func NewController(name string, description string) *Controller {
	return &Controller{
		Name:        name,
		Description: description,
	}
}

func (c *Controller) SetRoutes(route ...*Route) {
	c.Routes = append(c.Routes, route...)
}

//Проверяем на пустоту путь, если путь пуст то забираем из PkgPath
func (c *Controller) CheckPath(path string, v interface{}) string {
	if path == "" {
		path = c.getPkgPath(v)
	}
	return path
}

//Ищем все после пакета controllers
func (c *Controller) getPkgPath(v interface{}) string {
	t := reflect.TypeOf(v)
	pkg := strings.Replace(
		regexp.MustCompile(`controllers(.*)$`).FindString(t.PkgPath()),
		"controllers",
		"",
		-1,
	)
	return strings.Join([]string{pkg, strings.ToLower(t.Name())}, "/")
}
