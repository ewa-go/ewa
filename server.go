package egowebapi

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"time"
)

type Server struct {
	server      *fasthttp.Server
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

func New(name string, config Config) IServer {

	//Таймауты
	read, write, idle := config.Timeout.Get()
	//Инициализируем сервер
	server := &fasthttp.Server{
		Name:                               name,
		ReadTimeout:                        time.Duration(read) * time.Second,
		WriteTimeout:                       time.Duration(write) * time.Second,
		IdleTimeout:                        time.Duration(idle) * time.Second,
	}

	return &Server{
		Config: config,
		server: server,
	}
}

func (s *Server) Start() {
	go s.start()
}

func (s *Server) start() {

	//Устанавливаем статические файлы
	fs := fasthttp.FS{
		Root: s.Config.Root,
		GenerateIndexPages: true,
		Compress: true,
	}

	_ = fs.NewRequestHandler()

	//Устанавливаем роутер
	s.server.Handler = s.newRouter().Handler
	//Флаг старта
	s.Started = true
	//Если Secure == nil, то запускаем без сертификата
	if s.Config.Secure != nil {
		key, cert := s.Config.Secure.Get()
		if err := s.server.ListenAndServeTLS(fmt.Sprintf(":%d", s.Config.Port), cert, key); err != fasthttp.ErrConnectionClosed {
			s.server.Logger.Printf("%s", err)
		}
	} else {
		if err := s.server.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port)); err != fasthttp.ErrConnectionClosed {
			s.server.Logger.Printf("%s", err)
		}
	}
}

func (s *Server) newRouter() *router.Router {
	r := router.New()
	for _, c := range s.Controllers {
		for _, route := range c.Routes {
			r.Handle(route.Method, route.Path, route.Handler)
		}
	}
	return r
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
	return s.server.Shutdown()
}