package egowebapi

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber"
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
	//Controllers Controllers
}

type IServer interface {
	Start()
	Stop() error
	SetWeb(i IWeb, path string) *Server
	SetRest(i IRest, path string) *Server
}

func New(name string, config Config) (IServer, error) {

	//Таймауты
	read, write, idle := config.Timeout.Get()

	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	//Настройки
	settings := &fiber.Settings{
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

func (s *Server) Start() {
	go s.start()
}

func (s *Server) start() {

	//Флаг старта
	s.Started = true

	//Если Secure == nil, то запускаем без сертификата
	if s.Config.Secure != nil {
		//Формируем сертификат
		cert, err := tls.LoadX509KeyPair(s.Config.Secure.Get())
		if err != nil {
			//log
			return
		}
		//Запускаем слушатель с TLS настройкой
		if err := s.Listen(
			fmt.Sprintf(":%d", s.Config.Port),
			&tls.Config{Certificates: []tls.Certificate{cert}},
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

func (s *Server) rest(i IRest, method string, path string) {
	route := new(Route)
	method = strings.ToUpper(method)
	switch method {
	case "PUT":
		route = i.Put()
		break
	case "DELETE":
		route = i.Delete()
		break
	case "OPTIONS":
		route = i.Options()
		break
	default:
		s.web(i, method, path)
		return
	}

	if route != nil {
		for _, rpath := range route.Path {
			s.Add(method, p.Join(path, rpath), route.Handler)
		}
	}
}

func (s *Server) web(i IWeb, method string, path string) {
	route := new(Route)
	method = strings.ToUpper(method)
	switch method {
	case "GET":
		route = i.Get()
		break
	case "POST":
		route = i.Post()
		break
	}
	if route != nil {
		for _, rpath := range route.Path {
			s.Add(method, p.Join(path, rpath), route.Handler)
		}
	}
}

func (s *Server) SetWeb(i IWeb, path string) *Server {
	path = s.checkPath(path, i)
	s.web(i, "GET", path)
	s.web(i, "POST", path)

	return s
}

func (s *Server) SetRest(i IRest, path string) *Server {
	path = s.checkPath(path, i)
	s.SetWeb(i, path)
	s.rest(i, "PUT", path)
	s.rest(i, "DELETE", path)
	s.rest(i, "OPTIONS", path)

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
		regexp.MustCompile(`controllers(.*)$`).FindString(t.PkgPath()), //t.Elem().PkgPath()),
		"controllers",
		"",
		-1,
	)
	return strings.Join([]string{pkg, strings.ToLower(t.Name())}, "/")
}
