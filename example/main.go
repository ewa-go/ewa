package main

import (
	"errors"
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
	"github.com/egovorukhin/egowebapi/example/controllers/web/section1"
	__1 "github.com/egovorukhin/egowebapi/example/controllers/web/section1/1_1"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
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
	checkSession := func(key string) (string, string, error) {
		if value, ok := storage.GetStorage(key); ok {
			return value, "", nil
		}
		return "", "", errors.New("Элемент не найден")
	}
	//Обработчик ошибок
	errorHandler := func(ctx *fiber.Ctx, code int) error {
		text := "Unknown"
		if code == fiber.StatusForbidden {
			text = "Forbidden"
		}
		return ctx.Render("error", fiber.Map{"Code": code, "Text": text})
	}
	//Permission
	checkPermission := func(id interface{}, path string) bool {
		user, _ := storage.GetStorage(id.(string))
		if user == "user" && strings.Contains(path, "/section1/1_1") {
			return true
		}
		return false
	}

	//WEB
	cfg := ewa.Config{
		Port:    3005,
		Timeout: ewa.NewTimeout(30, 30, 30),
		Views: &ewa.Views{
			Directory: "views",
			Extension: ewa.Html,
			Engine: &ewa.Engine{
				Reload: false,
			},
		},
		Static: "views",
		Authorization: ewa.Authorization{
			Basic: &ewa.Basic{
				Handler: basicAuthHandler,
			},
		},
		Session: &ewa.Session{
			RedirectPath: "/login",
			Expires:      1 * time.Minute,
			Handler:      checkSession,
			ErrorHandler: errorHandler,
		},
		Permission: &ewa.Permission{
			Handler: checkPermission,
		},
	}
	// Указываем суффиксы
	suffix := ewa.NewSuffix(
		ewa.Suffix{Index: 2, Value: ":system"},
		ewa.Suffix{Index: 3, Value: ":version"},
	)
	//Инициализируем сервер
	ws, _ := ewa.New("Example", cfg)
	ws.Register(new(web.Home), "/")
	ws.Register(new(web.Login), "/login")
	ws.Register(new(web.Logout), "/logout")
	ws.Register(new(web.Swagger), "/info")
	ws.Register(new(__1.Document), "/section1/1_1/document")
	ws.Register(new(__1.List), "/section1/1_1/list")
	ws.Register(new(section1.Section_1_2), "/section1/1_2")
	ws.RegisterExt(new(api.User), "", "person", suffix...)
	ws.Register(new(api.WS), "")
	//webSocket
	ws.Register(new(controllers.WS), "")
	//wsserver.SetBasicAuth(ba)
	//Cors = nil - DefaultConfig
	ws.SetCors(nil)

	// Получаем объект fiber.App
	//wsserver.GetApp().Use()
	//wsserver.SetStore(nil)
	ws.Start()

	for {
		var input string
		_, err := fmt.Fscan(os.Stdin, &input)
		if err != nil {
			os.Exit(1)
		}
		switch strings.ToLower(input) {
		case "exit":
			fmt.Println(ws.Stop())
			os.Exit(0)
		}
	}
}
