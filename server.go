package egowebapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	p "path"
	"reflect"
	"regexp"
	"strings"
)

const (
	Name    = "EgoWebApi"
	Version = "v0.2.1"
)

//type Framework string

/*const (
	FrameworkFiber = "fiber"
	FrameworkEcho  = "echo"
)*/

type Server struct {
	Config    Config
	IsStarted bool
	webServer IServer
}

var swagger *Swagger

type IServer interface {
	Start(addr string, secure *Secure) error
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

//type Cors cors.Config

//type Store session.Config

/*type IServer interface {
	Start() error
	Stop() error
	Register(i interface{}, path string) *Server
	RegisterExt(i interface{}, path string, name string, suffix ...Suffix) *Server
	SetCors(config *Cors) *Server
	GetWebServer() IWebServer
	//SetStore(config *Store) * Server
}*/

func New(server IServer, config Config) *Server {

	//var server IServer
	//Таймауты
	//readTimeout, writeTimeout, idleTimeout := config.Timeout.Get()
	// Буферы
	//readBufferSize, writeBufferSize := config.BufferSize.Get()
	//Получаем расположение исполняемого файла
	/*exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}*/
	//Настройки
	/*settings := fiber.Config{
		BodyLimit:       config.BodyLimit,
		ReadTimeout:     time.Duration(readTimeout) * time.Second,
		WriteTimeout:    time.Duration(writeTimeout) * time.Second,
		IdleTimeout:     time.Duration(idleTimeout) * time.Second,
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
	}*/
	/*switch fw {
	case FrameworkFiber:
		// Указываем нужны ли страницы
		if config.Views != nil {
			if config.Views.Extension != None {
				settings.Views = config.Views.Extension.Engine( config.Views.Directory, config.Views.Engine)
			}
			if config.Views.Layout != "" {
				settings.ViewsLayout = config.Views.Layout
			}
		}
		//Инициализируем сервер
		server = &framework.Fiber{
			App: fiber.New(settings),
		}
	case FrameworkEcho:
		server = &framework.Echo{
			App: echo.New(),
		}
	}*/

	// Устанавливаем статические файлы
	if config.Static != nil {
		server.Static(config.Static.Prefix, config.Static.Root)
	}

	return &Server{
		Config:    config,
		webServer: server,
	}
}

// GetWebServer вернуть интерфейс веб сервера
func (s *Server) GetWebServer() interface{} {
	return s.webServer.GetApp()
}

// Start запуск сервера
func (s *Server) Start() error {
	//Флаг старта
	s.IsStarted = true
	// Получение адреса
	addr := fmt.Sprintf(":%d", s.Config.Port)
	// Запуск слушателя веб сервера
	return s.webServer.Start(addr, s.Config.Secure)
}

// Устанавливаем глобальные настройки для маршрутов
func (s *Server) newRoute() *Route {
	route := new(Route)
	if s.Config.Session != nil {
		route.isSession = s.Config.Session.AllRoutes
	}
	if s.Config.Permission != nil {
		route.isPermission = s.Config.Permission.AllRoutes
	}
	if s.Config.Authorization.AllRoutes != "" {
		route.auth = s.Config.Authorization.AllRoutes
	} else {
		route.auth = NoAuth
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

// Обрабатываем метод OPTIONS
/*func (s *Server) options(i IRestOptions, name, path string) {
	route := s.newRoute()
	i.Options(route)
	s.add(fiber.MethodOptions, name, path, route)
}

// Обрабатываем интерфейс IWeb
func (s *Server) web(i IWeb, name, path string) {
	s.get(i, name, path)
	s.post(i, name, path)
}

// Обрабатываем интерфейс IRest
func (s *Server) rest(i IRest, name, path string) {
	s.web(i, name, path)
	s.put(i, name, path)
	s.delete(i, name, path)
}

// Обрабатываем интерфейс IRestOptions
func (s *Server) restOptions(i IRestOptions, name, path string) {
	s.rest(i, name, path)
	s.options(i, name, path)
}*/

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
	//route.option.Method = method

	// Инициализируем Swagger
	if swagger == nil {
		http := "http"
		if s.Config.Secure != nil {
			http += "s"
		}
		addr := "127.0.0.1"
		swagger = &Swagger{
			Uri: fmt.Sprintf("%s://%s:%d", http, addr, s.Config.Port),
		}
	}
	// WebSocket
	/*if route.webSocket != nil && route.webSocket.UpgradeHandler != nil {
		s.webServer.Any(path, route.webSocket.UpgradeHandler)
	}*/

	// Получаем handler маршрута
	h := route.getHandler(s.Config, swagger)

	// Перебираем параметры адресной строки
	for _, param := range route.params {
		// Объединяем путь и параметры
		path = p.Join(path, param)
		// Добавляем метод, путь и обработчик
		s.webServer.Add(method, path, h)
		// Добавляем запись в swagger
		swagger.Add(name, path, route)
	}
}

// RegisterEx Регистрация интерфейсов
func (s *Server) RegisterEx(v interface{}, path string, name string, suffix ...Suffix) *Server {

	// Устанавливаем имя и путь
	name, path = s.getPkgNameAndPath(path, name, v, suffix...)
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

// SetCors Установка CORS
//TODO for fiber and Echo
/*func (s *Server) SetCors(config *Cors) *Server {
	cfg := cors.ConfigDefault
	if config != nil {
		cfg = cors.Config(*config)
	}
	s.webServer.Use(cors.New(cfg))
	return s
}*/

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
