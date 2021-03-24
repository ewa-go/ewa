package egowebapi

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	*fiber.App
	Name        string
	Started     bool
	Config      Config
	Controllers Controllers
}

type IServer interface {
	Start()
	Stop() error
	SetControllers(c Controllers) *Server
	GetControllers() Controllers
}

func New(name string, config Config) (IServer, error) {

	//Таймауты
	read, write, idle := config.Timeout.Get()

	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	//Инициализируем сервер
	server := fiber.New(&fiber.Settings{
		Views:        html.New(filepath.Join(filepath.Dir(exe), "views"), ".html"),
		ReadTimeout:  time.Duration(read) * time.Second,
		WriteTimeout: time.Duration(write) * time.Second,
		IdleTimeout:  time.Duration(idle) * time.Second,
	})

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

	//Устанавливаем статические файлы
	s.Static("/", "./static")

	//Устанавливаем роутер
	for _, c := range s.Controllers {
		for _, route := range c.Routes {
			s.Add(route.Method, route.Path, route.Handler)
		}
	}

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

func (s *Server) SetControllers(c Controllers) *Server {
	s.Controllers = append(s.Controllers, c...)
	return s
}

func (s *Server) GetControllers() Controllers {
	return s.Controllers
}

func (s *Server) Stop() error {
	s.Started = false
	return s.Shutdown()
}
