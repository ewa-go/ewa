package egowebapi

import (
	"fmt"
	"github.com/egovorukhin/egowebapi/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/websocket/v2"
	"github.com/valyala/fasthttp"
	"os"
	p "path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Server struct {
	*fiber.App
	Name      string
	IsStarted bool
	Config    Config
	Swagger   *swagger.Swagger
}

type Cors cors.Config
type Store session.Config

type IServer interface {
	Start()
	StartAsync()
	Stop() error
	Register(i interface{}, path string) *Server
	RegisterExt(i interface{}, path string, name string, suffix map[int]string) *Server
	SetCors(config *Cors) *Server
	GetApp() *fiber.App
	//SetStore(config *Store) * Server
}

func New(name string, config Config) (IServer, error) {

	//Таймауты
	read, write, idle := config.Timeout.Get()
	//Получаем расположение исполняемого файла
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	//Настройки
	settings := fiber.Config{
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
		/*ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			if config.Views != nil {
				err = ctx.Status(code).SendFile(fmt.Sprintf("/%d%s", code, config.Views.Extension))
				if err != nil {
					return ctx.Status(500).SendString("Internal Server Error")
				}
			} else {
				ctx.Status(code)
			}
			return nil
		},*/
	}
	//Указываем нужны ли страницы
	if config.Views != nil {
		if config.Views.Extension != None {
			settings.Views = config.Views.Extension.Engine(filepath.Join(filepath.Dir(exePath), config.Views.Directory), config.Views.Engine)
		}
		if config.Views.Layout != "" {
			settings.ViewsLayout = config.Views.Layout
		}
	}
	//Инициализируем сервер
	server := fiber.New(settings)
	//Устанавливаем статические файлы
	if config.Static != "" {
		server.Static("/", filepath.Join(filepath.Dir(exePath), config.Static))
	}

	return &Server{
		Name:   name,
		Config: config,
		App:    server,
	}, nil
}

func (s *Server) GetApp() *fiber.App {
	return s.App
}

func (s *Server) StartAsync() {
	go s.Start()
}

func (s *Server) Start() {

	//Флаг старта
	s.IsStarted = true

	//Если Secure == nil, то запускаем без сертификата
	if s.Config.Secure != nil {
		//Формируем сертификат
		cert, key := s.Config.Secure.Get()
		//Запускаем слушатель с TLS настройкой
		if err := s.ListenTLS(
			fmt.Sprintf(":%d", s.Config.Port),
			cert,
			key,
		); err != fasthttp.ErrConnectionClosed {
			//s.Logger.Printf("%s", err)
		}
		return
	}

	//Запускаем слушатель
	if err := s.Listen(fmt.Sprintf(":%d", s.Config.Port)); err != fasthttp.ErrConnectionClosed {
		//s.server.Logger.Printf("%s", err)
	}
}

func (s *Server) webSocket(i IWebSocket, path string) *swagger.Option {
	route := new(Route)
	i.Get(route)
	return s.add(fiber.MethodGet, path, route)
}

func (s *Server) rest(i IRest, method string, path string) *swagger.Option {
	route := new(Route)
	method = strings.ToUpper(method)
	switch method {
	case fiber.MethodPut:
		i.Put(route)
		break
	case fiber.MethodDelete:
		i.Delete(route)
		break
	default:
		s.web(i, method, path)
		return nil
	}
	return s.add(method, path, route)
}

func (s *Server) web(i IWeb, method string, path string) *swagger.Option {
	method = strings.ToUpper(method)
	route := new(Route)
	switch method {
	case fiber.MethodGet:
		i.Get(route)
		break
	case fiber.MethodPost:
		i.Post(route)
		break
	}
	return s.add(method, path, route)
}

func (s *Server) add(method string, path string, route *Route) *swagger.Option {

	// Если нет ни одного handler, то выходим
	if route.Handler == nil {
		return nil
	}
	/*if route.Handler == nil &&
		//route.WebHandler == nil &&
		route.LoginHandler == nil &&
		route.LogoutHandler == nil &&
		route.SwaggerHandler == nil &&
		(route.ws != nil && route.ws.Handler == nil) {
		return nil
	}*/

	// Инициализируем Swagger
	if s.Swagger == nil {
		s.Swagger = new(swagger.Swagger)
	}

	if route.Params == nil {
		route.Params = []string{""}
	}

	// Подключаем сессии
	_session := s.Config.Session

	// Получаем handler маршрута
	h := route.GetHandler(s.Config)

	for _, param := range route.Params {

		// Условно определяем что сессии и права на маршруты будут только для web страниц
		if route.Handler != nil {
			// Проверяем маршрут на актуальность сессии
			if (route.IsSession && _session != nil) || route.IsSession {
				h = _session.check(route.Handler, route.IsPermission)
			} else {
				h = route.Handler(ctx, nil)
			}
		}

		// Подключаем basic auth для api маршрутов
		/*if s.Config.Authorization.Basic != nil && route.Authorization
		if s.Config.BasicAuth != nil && route.IsBasicAuth {
			h = s.Config.BasicAuth.check(h)
		}*/
		// Авторизация - вход
		/*if _session != nil && route.LoginHandler != nil {
			h = _session.login(route.LoginHandler)
		}
		// Авторизация - выход
		if _session != nil && route.LogoutHandler != nil {
			h = _session.logout(route.LogoutHandler)
		}*/
		// WebSocket
		if route.ws != nil {
			if route.ws.UpgradeHandler != nil {
				s.Use(path, route.ws.UpgradeHandler)
			}
			if route.ws.Handler != nil {
				h = websocket.New(route.ws.Handler)
			}
		}

		// Заполняем Swagger
		if route.SwaggerHandler != nil {
			h = s.Swagger.check(route.SwaggerHandler)
		}

		s.Add(method, p.Join(path, param), h)
	}

	option := &swagger.Option{
		Params:      route.Params,
		Description: route.Description,
		Method:      method,
	}

	s.Swagger.AddOption(option)

	return option
}

// RegisterExt Регистрация интерфейсов
func (s *Server) RegisterExt(i interface{}, path string, name string, suffix map[int]string) *Server {
	// Проверка интерфейса на соответствие
	switch i.(type) {
	case IWebSocket:
		return s.registerWebSocket(i.(IWebSocket), path)
	case IWeb:
		return s.registerWeb(i.(IWeb), path)
	case IRest:
		return s.registerRest(i.(IRest), path, name, suffix)
	}
	return s
}

func (s *Server) Register(i interface{}, path string) *Server {
	return s.RegisterExt(i, path, "", nil)
}

// Регистрируем интерфейс IWebSocket
func (s *Server) registerWebSocket(i IWebSocket, path string) *Server {
	//Устанавливаем имя и путь
	_, path = s.getPkgNameAndPath(path, "", i, nil)
	s.webSocket(i, path)
	return s
}

// Регистрируем интерфейс IWeb
func (s *Server) registerWeb(i IWeb, path string) *Server {
	//Устанавливаем имя и путь
	_, path = s.getPkgNameAndPath(path, "", i, nil)

	s.web(i, fiber.MethodGet, path)
	s.web(i, fiber.MethodPost, path)

	return s
}

// Регистрируем интерфейс IRest
func (s *Server) registerRest(i IRest, path string, name string, suffix map[int]string) *Server {
	//Устанавливаем имя и путь
	name, path = s.getPkgNameAndPath(path, name, i, suffix)
	//Устанавливаем Swagger
	swagger := swagger.newSwagger(name, path)
	swagger.AddOption(s.web(i, fiber.MethodGet, path))
	swagger.AddOption(s.web(i, fiber.MethodPost, path))
	swagger.AddOption(s.rest(i, fiber.MethodPut, path))
	swagger.AddOption(s.rest(i, fiber.MethodDelete, path))
	// Создаем исполнителя для метода Options
	s.Add(fiber.MethodOptions, path, i.Options(swagger))

	return s
}

// SetCors Установка CORS
func (s *Server) SetCors(config *Cors) *Server {
	cfg := cors.ConfigDefault
	if config != nil {
		cfg = cors.Config(*config)
	}
	s.Use(cors.New(cfg))
	return s
}

// Stop Остановка сервера
func (s *Server) Stop() error {
	s.IsStarted = false
	return s.Shutdown()
}

//Ищем все после пакета controllers
func (s *Server) getPkgNameAndPath(path, name string, v interface{}, suffix map[int]string) (string, string) {
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
		for index, val := range suffix {
			array = s.insert(array, index, val)
		}
		path = strings.Join(array, "/") + "/" + strings.ToLower(name)
	}

	return strings.Title(name), path
}

func (s *Server) insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
