package main

import (
	"errors"
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/fiber/controllers/web"
	"github.com/egovorukhin/egowebapi/example/fiber/src/storage"
	f "github.com/egovorukhin/egowebapi/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//BasicAuth
	basicAuthHandler := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}
	//Session
	checkSession := func(key string) (string, error) {
		if value, ok := storage.GetStorage(key); ok {
			return value, nil
		}
		return "", errors.New("Элемент не найден")
	}
	//Обработчик ошибок
	errorHandler := func(c *ewa.Context, code int, err interface{}) error {
		return c.Render("error", map[string]interface{}{"Code": code, "Text": err})
	}
	//Permission
	checkPermission := func(username, path string) bool {
		if username == "user" {
			switch path {
			case "/":
				return true
			}
		}
		return false
	}

	root := "./views"

	// Fiber
	app := fiber.New(fiber.Config{
		Views: f.NewViews(root, f.Html, &f.Engine{
			Reload: true,
		}),
	})
	app.Use(favicon.New(favicon.Config{
		File: "./dist/favicon.ico",
	}))
	app.Use(cors.New())
	server := &f.Server{App: app}
	// Конфиг
	cfg := ewa.Config{
		Port: 3005,
		Static: &ewa.Static{
			Prefix: "/",
			Root:   root,
		},
		Views: &ewa.Views{
			Root:   root,
			Engine: f.Html,
		},
		Authorization: ewa.Authorization{
			Basic: basicAuthHandler,
		},
		Session: &ewa.Session{
			RedirectPath:   "/login",
			Expires:        1 * time.Minute,
			SessionHandler: checkSession,
			ErrorHandler:   errorHandler,
		},
		Permission: &ewa.Permission{
			Handler: checkPermission,
		},
	}
	// Указываем суффиксы
	/*suffix := ewa.NewSuffix(
		ewa.Suffix{Index: 2, Value: ":system"},
		ewa.Suffix{Index: 3, Value: ":version"},
	)*/
	//Инициализируем сервер
	ws := ewa.New(server, cfg)
	ws.Register(new(web.Index), "/")
	ws.Register(new(web.Home), "/home")
	ws.Register(new(web.Login), "/login")
	ws.Register(new(web.Logout), "/logout")
	//ws.RegisterEx(new(api2.Username), "", "person", suffix...)
	//ws.Register(new(api2.WS), "")
	//webSocket
	//ws.Register(new(controllers.WS), "")

	// Канал для получения ошибки, если таковая будет
	errChan := make(chan error, 2)
	go func() {
		errChan <- ws.Start()
	}()

	// Ждем сигнал от ОС
	go getSignal(errChan)

	fmt.Println("Старт приложения")
	fmt.Printf("Остановка приложения. %s", <-errChan)

}

func getSignal(errChan chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	errChan <- fmt.Errorf("%s", <-c)
}
