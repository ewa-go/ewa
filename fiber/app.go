package fiber

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Server struct {
	App *fiber.App
}

func (s *Server) Start(addr string, secure *ewa.Secure) (err error) {

	// Если флаг для безопасности true, то запускаем механизм с TLS
	if secure != nil {
		// Возвращаем данные по сертификату
		cert, key := secure.Get()
		// Запускаем слушатель с TLS настройкой
		err = s.App.ListenTLS(addr, cert, key)
	} else {
		err = s.App.Listen(addr)
	}
	if err != nil && err != fasthttp.ErrConnectionClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.App.Shutdown()
}

func (s *Server) Static(prefix, root string) {
	s.App.Static(prefix, root)
}

func (s *Server) Any(path string, handler interface{}) {
	if h, ok := handler.(fiber.Handler); ok {
		s.App.Use(path, h)
	}
}

func (s *Server) Use(params ...interface{}) {
	s.App.Use(params)
}

func (s *Server) Add(method, path string, handler ewa.Handler) {
	s.App.Add(method, path, func(ctx *fiber.Ctx) error {
		c := ewa.NewContext(&Context{Ctx: ctx})
		return handler(c)
	})
}

func (s *Server) GetApp() interface{} {
	return s.App
}

func (s *Server) NotFoundPage(path, page string) {
	s.App.Use(func(ctx *fiber.Ctx) error {
		return ctx.Render(page, nil)
	})
}
