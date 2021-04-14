package egowebapi

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
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
	BasicAuth
}

type IServer interface {
	Start()
	StartAsync()
	Stop() error
	RegisterWeb(i IWeb, path string) *Server
	RegisterRest(i IRest, path string, name string) *Server
	SetBasicAuth(auth BasicAuth) *Server
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
	} else {
		//Запускаем слушатель
		if err := s.Listen(fmt.Sprintf(":%d", s.Config.Port)); err != fasthttp.ErrConnectionClosed {
			//s.server.Logger.Printf("%s", err)
		}
	}
}

func (s *Server) SetBasicAuth(auth BasicAuth) *Server {
	s.Use(basicauth.New(auth.Config))
	return s
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
	if route == nil {
		return nil
	}

	if route.Params == nil {
		route.Params = []string{""}
	}

	for _, rpath := range route.Params {
		s.Add(method, p.Join(path, rpath), route.Handler)
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

func (s *Server) RegisterRest(i IRest, path string, name string) *Server {
	//Устанавливаем имя и путь
	name, path = s.getPkgNameAndPath(path, name, i)
	//Устанавливаем Swagger
	swagger := newSwagger(name, path)
	swagger.AddOption(s.web(i, fiber.MethodGet, path))
	swagger.AddOption(s.web(i, fiber.MethodPost, path))
	swagger.AddOption(s.rest(i, fiber.MethodPut, path))
	swagger.AddOption(s.rest(i, fiber.MethodDelete, path))
	//Создаем исполнитеоля для Options
	s.Add(fiber.MethodOptions, path, i.Options(swagger))

	return s
}

func (s *Server) Stop() error {
	s.Started = false
	return s.Shutdown()
}

//Ищем все после пакета controllers
func (s *Server) getPkgNameAndPath(path, name string, v interface{}) (string, string) {
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
	if path == "" {
		path = strings.Join([]string{pkg, strings.ToLower(t.Name())}, "/")
	}
	if name == "" {
		name = t.Name()
	}
	return name, path
}
