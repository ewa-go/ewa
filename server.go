package egowebapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
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
	Name    string
	Started bool
	Config  Config
	//store   *session.Store
}

type Cors cors.Config
type Store session.Config

type IServer interface {
	Start()
	StartAsync()
	Stop() error
	RegisterWeb(i IWeb, path string) *Server
	RegisterRest(i IRest, path string, name string, suffix ...Suffix) *Server
	SetCors(config *Cors) *Server
	//SetStore(config *Store) * Server
}

func New(name string, config Config) (IServer, error) {

	//Таймауты
	read, write, idle := config.Timeout.Get()
	//Получаем расположение исполняемого файла
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	//Настройки
	settings := fiber.Config{
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
	}
	//Указываем нужны ли страницы
	if config.Views != nil {
		settings.Views = html.New(filepath.Join(filepath.Dir(exe), config.Views.Root), config.Views.Engine)
	}
	//Инициализируем сервер
	server := fiber.New(settings)
	//Устанавливаем статические файлы
	if config.Static != "" {
		server.Static("/", filepath.Join(filepath.Dir(exe), config.Static))
	}

	return &Server{
		Name:   name,
		Config: config,
		App:    server,
	}, nil
}

func (s *Server) StartAsync() {
	go s.Start()
}

func (s *Server) Start() {

	//Флаг старта
	s.Started = true

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

func (s *Server) rest(i IRest, method string, path string) *Option {
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

func (s *Server) web(i IWeb, method string, path string) *Option {
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

func (s *Server) add(method string, path string, route *Route) *Option {
	if route.Handler == nil && route.WebHandler == nil {
		return nil
	}

	if route.Params == nil {
		route.Params = []string{""}
	}

	for _, param := range route.Params {
		h := route.Handler
		//Подключаем basic auth
		if s.Config.BasicAuth != nil && route.IsBasicAuth {
			h = s.Config.BasicAuth.check(h)
		}
		//Подключаем сессии
		if route.IsSession && route.WebHandler != nil {
			h = s.Config.Session.check(route.WebHandler)
		}
		//Проверка разрешений
		if route.IsPermission {
			h = s.Config.Permission.check(h)
		}

		s.Add(method, p.Join(path, param), h)
	}

	return &Option{
		Params:      route.Params,
		Description: route.Description,
		Method:      method,
	}
}

func (s *Server) RegisterWeb(i IWeb, path string) *Server {
	//Устанавливаем имя и путь
	_, path = s.getPkgNameAndPath(path, "", i)

	s.web(i, fiber.MethodGet, path)
	s.web(i, fiber.MethodPost, path)

	return s
}

type Suffix struct {
	Index int
	Value string
}

func (s *Server) RegisterRest(i IRest, path string, name string, suffix ...Suffix) *Server {
	//Устанавливаем имя и путь
	name, path = s.getPkgNameAndPath(path, name, i, suffix...)
	//Устанавливаем Swagger
	swagger := newSwagger(name, path)
	swagger.AddOption(s.web(i, fiber.MethodGet, path))
	swagger.AddOption(s.web(i, fiber.MethodPost, path))
	swagger.AddOption(s.rest(i, fiber.MethodPut, path))
	swagger.AddOption(s.rest(i, fiber.MethodDelete, path))
	//Создаем исполнителя для Options
	s.Add(fiber.MethodOptions, path, i.Options(swagger))

	return s
}

func (s *Server) SetCors(config *Cors) *Server {
	cfg := cors.ConfigDefault
	if config != nil {
		cfg = cors.Config(*config)
	}
	s.Use(cors.New(cfg))
	return s
}

/*func (s *Server) SetStore(config *Store) *Server {
	cfg := session.ConfigDefault
	if config != nil {
		cfg = session.Config(*config)
	}
	s.store = session.New(cfg)
	return s
}*/

func (s *Server) Stop() error {
	s.Started = false
	return s.Shutdown()
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
		for _, s := range suffix {
			array = insert(array, s.Index, s.Value)
		}
		path = strings.Join(array, "/") + "/" + strings.ToLower(name)
	}

	return strings.Title(name), path
}

func insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
