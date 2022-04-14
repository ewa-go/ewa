package egowebapi

import (
	"fmt"
	"github.com/alecthomas/jsonschema"
	"github.com/egovorukhin/egowebapi/security"
	p "path"
	"regexp"
	"strings"
)

const (
	Name    = "EgoWebApi"
	Version = "v0.2.6"
)

type Server struct {
	Config      Config
	IsStarted   bool
	WebServer   IServer
	Controllers []*Controller
	Swagger     *Swagger
}

type IServer interface {
	Start(addr string) error
	StartTLS(addr, cert, key string) error
	Stop() error
	Static(prefix, root string)
	Any(path string, handler interface{})
	Use(params ...interface{})
	Add(method, path string, handler Handler)
	GetApp() interface{}
	NotFoundPage(path, page string)
	ConvertParam(param string) string
}

type Suffix struct {
	Index int
	Value string
}

func NewSuffix(suffix ...Suffix) (s []Suffix) {
	for _, item := range suffix {
		s = append(s, item)
	}
	return
}

func New(server IServer, config Config) *Server {

	// Устанавливаем статические файлы
	if config.Static != nil {
		server.Static(config.Static.Prefix, config.Static.Root)
	}

	s := &Server{
		Config:    config,
		WebServer: server,
		Swagger: &Swagger{
			Swagger:             "2.0",
			Host:                fmt.Sprintf("localhost:%d", config.Port),
			BasePath:            "/",
			SecurityDefinitions: SecurityDefinitions{},
			Paths:               Paths{},
			Definitions:         map[string]*jsonschema.Type{},
		},
	}

	return s
}

// GetWebServer вернуть интерфейс веб сервера
func (s *Server) GetWebServer() interface{} {
	return s.WebServer.GetApp()
}

// Start запуск сервера
func (s *Server) Start() (err error) {

	for _, v := range s.Controllers {

		v.initialize()
		path := v.Path
		name := v.Tag.Name

		// Добавляем тэги контроллера
		s.Swagger.Tags = append(s.Swagger.Tags, v.Tag)

		// Проверка интерфейса на соответствие
		if i, ok := v.Interface.(IGet); ok {
			err = s.get(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IPost); ok {
			err = s.post(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IPut); ok {
			err = s.put(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IDelete); ok {
			err = s.delete(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IOptions); ok {
			err = s.options(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IPatch); ok {
			err = s.patch(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IHead); ok {
			err = s.head(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(IConnect); ok {
			err = s.connect(i, name, path)
			if err != nil {
				return
			}
		}
		if i, ok := v.Interface.(ITrace); ok {
			err = s.trace(i, name, path)
			if err != nil {
				return
			}
		}
	}

	// Схема
	scheme := "http"
	//Флаг старта
	s.IsStarted = true
	// Получение адреса
	addr := fmt.Sprintf(":%d", s.Config.Port)
	// Если флаг для безопасности true, то запускаем механизм с TLS
	if s.Config.Secure != nil {
		// Security
		scheme += "s"
		// Возвращаем данные по сертификату
		cert, key := s.Config.Secure.Get()
		// Запускаем слушатель с TLS настройкой
		return s.WebServer.StartTLS(addr, cert, key)
	}
	// Добавляем схему в Swagger
	if s.Swagger != nil {
		s.Swagger.SetSchemes(scheme)
	}
	// Запуск слушателя веб сервера
	return s.WebServer.Start(addr)
}

// Stop Остановка сервера
func (s *Server) Stop() error {
	s.IsStarted = false
	return s.WebServer.Stop()
}

// Устанавливаем глобальные настройки для маршрутов
func (s *Server) newRoute() *Route {

	route := &Route{
		Operation: Operation{
			Responses: map[string]Response{
				"default": {
					Description: "successful operation",
				},
			},
		},
	}
	if s.Config.Permission != nil {
		route.isPermission = s.Config.Permission.AllRoutes
	}
	if s.Config.Authorization.AllRoutes != security.NoAuth {
		route.SetSecurity(s.Config.Authorization.AllRoutes)
	}

	return route
}

// Обрабатываем метод GET
func (s *Server) get(i IGet, name, path string) error {
	route := s.newRoute()
	i.Get(route)
	return s.add(MethodGet, name, path, route)
}

// Обрабатываем метод POST
func (s *Server) post(i IPost, name, path string) error {
	route := s.newRoute()
	i.Post(route)
	return s.add(MethodPost, name, path, route)
}

// Обрабатываем метод PUT
func (s *Server) put(i IPut, name, path string) error {
	route := s.newRoute()

	i.Put(route)
	return s.add(MethodPut, name, path, route)
}

// Обрабатываем метод DELETE
func (s *Server) delete(i IDelete, name, path string) error {
	route := s.newRoute()
	i.Delete(route)
	return s.add(MethodDelete, name, path, route)
}

// Обрабатываем метод OPTIONS
func (s *Server) options(i IOptions, name, path string) error {
	route := s.newRoute()
	i.Options(route)
	return s.add(MethodOptions, name, path, route)
}

// Обрабатываем метод PATCH
func (s *Server) patch(i IPatch, name, path string) error {
	route := s.newRoute()
	i.Patch(route)
	return s.add(MethodPatch, name, path, route)
}

// Обрабатываем метод HEAD
func (s *Server) head(i IHead, name, path string) error {
	route := s.newRoute()
	i.Head(route)
	return s.add(MethodHead, name, path, route)
}

// Обрабатываем метод CONNECT
func (s *Server) connect(i IConnect, name, path string) error {
	route := s.newRoute()
	i.Connect(route)
	return s.add(MethodConnect, name, path, route)
}

// Обрабатываем метод TRACE
func (s *Server) trace(i ITrace, name, path string) error {
	route := s.newRoute()
	i.Trace(route)
	return s.add(MethodTrace, name, path, route)
}

// Добавить маршрут в веб сервер
func (s *Server) add(method, tagName, path string, route *Route) error {

	// Если нет ни одного handler, то выходим
	if route.Handler == nil {
		return nil
	}

	params := route.Operation.GetParams()

	if params == nil || route.isEmptyParam {
		params = append(params, "")
	}

	/*var view *View
	// Проверка на view
	if s.Config.Views != nil {
		files, _ := filepath.Glob(filepath.Join(s.Config.Views.Root, strings.ToLower(name)+s.Config.Views.Engine))
		for _, file := range files {
			view = &View{
				Filename: strings.Replace(filepath.Base(file), s.Config.Views.Engine, "", -1),
				Filepath: file,
				Layout:   s.Config.Views.Layout,
			}
		}
	}*/

	// Авторизация в swagger
	for _, sec := range route.Security {
		for key := range sec {
			s.Swagger.setSecurityDefinition(key, s.Config.Authorization.Get(key).Definition())
		}
	}

	// Добавляем ссылку на тэг в контроллере
	route.Operation.Tags = append(route.Operation.Tags, tagName)

	// Получаем handler маршрута
	h := route.getHandler(s.Config, nil, *s.Swagger)

	// Перебираем параметры адресной строки
	for _, param := range params {

		// Объединяем путь и параметры
		fullPath := p.Join(path, param)

		// Добавляем пути и методы в swagger
		s.Swagger.setPath(fullPath, strings.ToLower(method), route.Operation)

		// Проверка на пустые пути
		if param != "" {
			matches := regexp.MustCompile(`{(\w+)}`).FindStringSubmatch(fullPath)
			if len(matches) == 2 {
				fullPath = strings.ReplaceAll(fullPath, matches[0], s.WebServer.ConvertParam(matches[1]))
			}
		}

		// Добавляем метод, путь и обработчик
		s.WebServer.Add(method, fullPath, h)
	}

	return nil
}

// Register Регмтсрация контроллера
func (s *Server) Register(i interface{}) *Controller {
	controller := &Controller{
		Interface: i,
	}
	s.Controllers = append(s.Controllers, controller)
	return controller
}

// Функция вернет Имя и версию
func (s *Server) String() string {
	return fmt.Sprintf("%s %s", Name, Version)
}
