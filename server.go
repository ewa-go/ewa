package egowebapi

import (
	"fmt"
	v2 "github.com/egovorukhin/egowebapi/swagger/v2"
	"github.com/gofiber/fiber/v2"
	p "path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const (
	Name    = "EgoWebApi"
	Version = "v0.2.5"
)

type Server struct {
	Config    Config
	IsStarted bool
	webServer IServer
	Swagger   *v2.Swagger
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
		webServer: server,
	}

	// Инициализация swagger
	if config.Swagger != nil {
		s.Swagger = v2.New(config.Swagger.Host, config.Swagger.Info)
	}

	return s
}

// GetWebServer вернуть интерфейс веб сервера
func (s *Server) GetWebServer() interface{} {
	return s.webServer.GetApp()
}

// Start запуск сервера
func (s *Server) Start() (err error) {
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
		return s.webServer.StartTLS(addr, cert, key)
	}
	// Добавляем схему в Swagger
	if s.Swagger != nil {
		s.Swagger.Schemes = append(s.Swagger.Schemes, scheme)
	}
	// Запуск слушателя веб сервера
	return s.webServer.Start(addr)
}

// Устанавливаем глобальные настройки для маршрутов
func (s *Server) newRoute() *Route {
	route := new(Route)
	if s.Config.Permission != nil {
		route.isPermission = s.Config.Permission.AllRoutes
	}
	if s.Config.Authorization.AllRoutes != "" {
		route.auth = append(route.auth, s.Config.Authorization.AllRoutes)
	}

	return route
}

// Обрабатываем метод GET
func (s *Server) get(i IGet, name, path string) {
	route := s.newRoute()
	i.Get(route)
	s.add(fiber.MethodGet, name, path, route)
}

// Обрабатываем метод POST
func (s *Server) post(i IPost, name, path string) {
	route := s.newRoute()
	i.Post(route)
	s.add(fiber.MethodPost, name, path, route)
}

// Обрабатываем метод PUT
func (s *Server) put(i IPut, name, path string) {
	route := s.newRoute()
	i.Put(route)
	s.add(fiber.MethodPut, name, path, route)
}

// Обрабатываем метод DELETE
func (s *Server) delete(i IDelete, name, path string) {
	route := s.newRoute()
	i.Delete(route)
	s.add(fiber.MethodDelete, name, path, route)
}

// Обрабатываем метод OPTIONS
func (s *Server) options(i IOptions, name, path string) {
	route := s.newRoute()
	i.Options(route)
	s.add(fiber.MethodOptions, name, path, route)
}

// Обрабатываем метод PATCH
func (s *Server) patch(i IPatch, name, path string) {
	route := s.newRoute()
	i.Patch(route)
	s.add(fiber.MethodPatch, name, path, route)
}

// Обрабатываем метод HEAD
func (s *Server) head(i IHead, name, path string) {
	route := s.newRoute()
	i.Head(route)
	s.add(fiber.MethodHead, name, path, route)
}

// Обрабатываем метод CONNECT
func (s *Server) connect(i IConnect, name, path string) {
	route := s.newRoute()
	i.Connect(route)
	s.add(fiber.MethodConnect, name, path, route)
}

// Обрабатываем метод TRACE
func (s *Server) trace(i ITrace, name, path string) {
	route := s.newRoute()
	i.Trace(route)
	s.add(fiber.MethodTrace, name, path, route)
}

func (s *Server) add(method string, name, path string, route *Route) {

	// Если нет ни одного handler, то выходим
	if route == nil || route.Handler == nil || method == "" {
		return
	}

	if route.params == nil {
		route.params = []string{"", "/"}
	} else {
		// Проверка пути на пустоту и слэш
		emptyPath := false
		slash := false
		for _, param := range route.params {
			switch param {
			case "":
				emptyPath = true
				break
			case "/":
				slash = true
				break
			}
		}
		if emptyPath && !slash {
			route.params = append(route.params, "/")
		} else if !emptyPath && slash {
			route.params = append(route.params, "")
		}
	}

	var view *View
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
	}

	// Добавляем запись в swagger
	if s.Swagger != nil {
		// Модели
		for _, param := range route.path.Parameters {
			if param.Type == v2.ParameterTypeBody {
				if param.Body != nil {
					Type := ""
					switch reflect.TypeOf(param.Body).Kind() {
					case reflect.Slice, reflect.Array:
						Type = "array"
					}
					param.Schema = v2.Schema{
						Type: Type,
						Items: map[string]string{
							"$ref": "#/definitions/" + name,
						},
					}
				}
			}
		}
		// Авторизация
		for _, a := range route.auth {
			sec, sd := s.GetSecurityDefinitions(a)
			route.path.Security = sec
			s.Swagger.SecurityDefinitions = sd
		}
		// Добавляем тэги, контролеры
		route.path.Tags = append(route.path.Tags, name)
		// Добавляем методы
		if s.Swagger.Paths[path] == nil {
			s.Swagger.Paths[path] = v2.Methods{}
		}
		s.Swagger.Paths[path][strings.ToLower(method)] = route.path
	}

	// Получаем handler маршрута
	h := route.getHandler(s.Config, view, s.Swagger)

	// Перебираем параметры адресной строки
	for _, param := range route.params {
		// Объединяем путь и параметры
		path = p.Join(path, param)
		// Добавляем метод, путь и обработчик
		s.webServer.Add(method, path, h)
	}
}

// RegisterEx Регистрация интерфейсов
func (s *Server) RegisterEx(v interface{}, path string, name string, suffix ...Suffix) *Server {

	// Устанавливаем имя и путь
	name, path = s.getPkgNameAndPath(path, name, v, suffix...)
	// Заполняем Tag для Swagger
	if s.Swagger != nil {
		if i, ok := v.(ITag); ok {
			tag := i.Tag()
			if tag.Name == "" {
				tag.Name = name
			}
			s.Swagger.Tags = append(s.Swagger.Tags, tag)
		}
	}
	// Проверка интерфейса на соответствие
	if i, ok := v.(IGet); ok {
		s.get(i, name, path)
	}
	if i, ok := v.(IPost); ok {
		s.post(i, name, path)
	}
	if i, ok := v.(IPut); ok {
		s.put(i, name, path)
	}
	if i, ok := v.(IDelete); ok {
		s.delete(i, name, path)
	}
	if i, ok := v.(IOptions); ok {
		s.options(i, name, path)
	}
	if i, ok := v.(IPatch); ok {
		s.patch(i, name, path)
	}
	if i, ok := v.(IHead); ok {
		s.head(i, name, path)
	}
	if i, ok := v.(IConnect); ok {
		s.connect(i, name, path)
	}
	if i, ok := v.(ITrace); ok {
		s.trace(i, name, path)
	}

	// Страница 404
	// TODO NotFound
	if s.Config.NotFoundPage != "" {
		s.webServer.NotFoundPage(path, s.Config.NotFoundPage)
	}

	return s
}

func (s *Server) Register(i interface{}, path string) *Server {
	return s.RegisterEx(i, path, "")
}

// Stop Остановка сервера
func (s *Server) Stop() error {
	s.IsStarted = false
	return s.webServer.Stop()
}

//Ищем все после пакета controllers
func (s *Server) getPkgNameAndPath(path, name string, v interface{}, suffix ...Suffix) (string, string) {
	//Если имя и путь установлены вручную, то выходим
	if path != "" && name != "" {
		return name, path
	}
	//Извлекаем имя и путь до controllers
	var t reflect.Type
	value := reflect.ValueOf(v)
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
	if name == "" {
		name = strings.Title(t.Name())
	}

	if path == "" {
		array := strings.Split(pkg, "/")
		for _, item := range suffix {
			array = s.insert(array, item.Index, item.Value)
		}
		path = strings.Join(array, "/") + "/" + strings.ToLower(name)
	}

	return strings.Title(name), path
}

func (s *Server) insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	} else if len(a) < index {
		return a
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func (s *Server) String() string {
	return fmt.Sprintf("%s %s", Name, Version)
}

func (s *Server) GetSecurityDefinitions(auth string) (sec v2.Security, sd v2.SecurityDefinitions) {

	sd = v2.SecurityDefinitions{}
	secure := v2.Secure{}
	a := s.Config.Authorization
	switch auth {
	case BasicAuth:
		sd[BasicAuth] = v2.SecurityDefinition{
			Type:        "basic",
			Description: "Basic Authorization",
		}
		secure[BasicAuth] = []string{}
		break
	case ApiKeyAuth:
		if a.ApiKey != nil {
			name, param := a.ApiKey.Get()
			sd[ApiKeyAuth] = v2.SecurityDefinition{
				Type:        "apiKey",
				Description: "Api Key Authorization",
				Name:        name,
				In:          param,
			}
			secure[ApiKeyAuth] = []string{}
		}
		break
		//TODO OAuth2 check
	case OAuth2Auth:
		values := []string{"write:pets", "read:pets"}
		secure[ApiKeyAuth] = values
		sd[OAuth2Auth] = v2.SecurityDefinition{
			Type:             "oauth2",
			Description:      "OAuth2 Authorization",
			Flow:             "oauth2",
			AuthorizationUrl: "",
			TokenUrl:         "",
			Scopes:           nil,
		}
	}

	sec = append(sec, secure)

	return sec, sd
}
