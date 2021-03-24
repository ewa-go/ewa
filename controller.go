package egowebapi

import (
	"html/template"
	"io"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type Controller struct {
	Name string
	Description string
	Routes Routes
}

type Controllers []*Controller

func NewController(name string, description string) *Controller {
	return &Controller{
		Name:        name,
		Description: description,
	}
}

//Возвращаем html страницу
//Используется для страниц Views, рендеринг страниц
func (c *Controller) View(wr io.Writer, pageName string, data interface{}) error {
	if pageName == "" {
		pageName = c.checkPath("", c) + ".html"
	}

	tmpl, err := template.ParseFiles(
		path.Join("views/share", "layout.html"),
		path.Join("views", pageName))
	if err != nil {
		return c.View(wr, "share/error.html", err.Error())
	}

	err = tmpl.ExecuteTemplate(wr, "layout", data)
	if err != nil {
		return err
	}
	return nil
}

//Отдаём страницы которые находяться в папке www.
//Т.е. используя функцию Page мы из модели MVC убираем views, а так же static
//в том виде который используется для MVC. Такой подход был реализован для проектов на React JS, Vue JS.
//Собираем проект React App с помощью npm или yarn, копируем все содержимое
//каталога build в каталог www вашего проекта и все будет работать. не забудьте создать
//контроллер Index, и добавить все возможные пути (react-router-dom, vue-router) это очень важно.
func (c *Controller) Page(wr io.Writer, pageName string, data interface{}) error {
	if pageName == "" {
		pageName = c.checkPath("", c)
	}
	pageName += ".html"

	page, err := template.ParseFiles(path.Join("www", pageName))
	if err != nil {
		return err
	}
	return page.Execute(wr, data)
}

//Проверяем на пустоту путь, если путь пуст то забираем из PkgPath
func (c * Controller) checkPath(path string, v interface{}) string {
	if path == "" {
		path = c.getPkgPath(v)
	}
	return path
}

//Ищем все после пакета controllers
func (c * Controller) getPkgPath(v interface{}) string {
	t := reflect.TypeOf(v)
	pkg := strings.Replace(
		regexp.MustCompile(`controllers(.*)$`).FindString(t.PkgPath()),
		"controllers",
		"",
		-1,
	)
	return strings.Join([]string{pkg, strings.ToLower(t.Name())}, "/")
}
