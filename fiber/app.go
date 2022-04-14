package fiber

import (
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
}

func (s *Server) Start(addr string) (err error) {
	return s.App.Listen(addr)
}

func (s *Server) StartTLS(addr, cert, key string) (err error) {
	// Запускаем слушатель с TLS настройкой
	return s.App.ListenTLS(addr, cert, key)
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
		return handler(ewa.NewContext(&Context{Ctx: ctx}))
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

func (s *Server) ConvertParam(param string) string {
	return ":" + param
}
