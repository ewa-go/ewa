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
	//Controllers Controllers
}

type IServer interface {
	Start()
	StartAsync()
	Stop() error
	RegisterWeb(i IWeb, path string) *Server
	RegisterRest(i IRest, path string) *Server
	SetBasicAuth(auth BasicAuth) *Server
}

func New(name string, config Config) (IServer, error) {

	//Таймауты
	read, write, idle := config.Timeout.Get()

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
		settings.Views = html.New(filepath.Join(filepath.Dir(exe), config.Views.Root), config.Views.Ext)
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

func (s *Server) rest(i IRest, method string, path string) {
	route := new(Route)
	method = strings.ToUpper(method)
	switch method {
	case fiber.MethodPut:
		route = i.Put()
		break
	case fiber.MethodDelete:
		route = i.Delete()
		break
	case fiber.MethodOptions:
		route = i.Options()
		break
	default:
		s.web(i, method, path)
		return
	}
	s.add(method, path, route)
}

func (s *Server) web(i IWeb, method string, path string) {
	route := new(Route)
	method = strings.ToUpper(method)
	switch method {
	case fiber.MethodGet:
		route = i.Get()
		break
	case fiber.MethodPost:
		route = i.Post()
		break
	}
	s.add(method, path, route)
}

func (s *Server) add(method string, path string, route *Route) {
	if route == nil {
		return
	}

	if route.Path == nil {
		route.Path = []string{""}
	}

	for _, rpath := range route.Path {
		s.Add(method, p.Join(path, rpath), route.Handler)
	}
}

func (s *Server) RegisterWeb(i IWeb, path string) *Server {
	path = s.checkPath(path, i)
	s.web(i, fiber.MethodGet, path)
	s.web(i, fiber.MethodPost, path)

	return s
}

func (s *Server) RegisterRest(i IRest, path string) *Server {
	path = s.checkPath(path, i)
	s.RegisterWeb(i, path)
	s.rest(i, fiber.MethodPut, path)
	s.rest(i, fiber.MethodDelete, path)
	s.rest(i, fiber.MethodOptions, path)

	return s
}

func (s *Server) Stop() error {
	s.Started = false
	return s.Shutdown()
}

//Проверяем на пустоту путь, если путь пуст то забираем из PkgPath
func (s *Server) checkPath(path string, v interface{}) string {
	if path == "" {
		path = s.getPkgPath(v)
	}
	return path
}

//Ищем все после пакета controllers
func (s *Server) getPkgPath(v interface{}) string {
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
	return strings.Join([]string{pkg, strings.ToLower(t.Name())}, "/")
}
