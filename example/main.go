package main

import (
	"errors"
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
	"github.com/egovorukhin/egowebapi/example/controllers/web/section1"
	"github.com/egovorukhin/egowebapi/example/src/storage"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

func main() {

	//BasicAuth
	authorizer := func(user string, pass string) bool {
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
	errorHandler := func(ctx *fiber.Ctx, code int, err string) error {
		return ctx.Render("error", fiber.Map{"Code": code, "Text": err})
	}
	//Permission
	checkPermission := func(route string) bool {
		if route != "" {
			return false
		}
		return true
	}
	//WEB
	cfg := ewa.Config{
		Port:    3005,
		Timeout: ewa.NewTimeout(30, 30, 30),
		Views: &ewa.Views{
			Root:   "views",
			Engine: ".html",
		},
		Static:    "views",
		BasicAuth: ewa.NewBasicAuth(authorizer, nil),
		Session: &ewa.Session{
			RedirectPath: "/login",
			Check:        checkSession,
		},
		Permission: &ewa.Permission{
			Check: checkPermission,
			Error: errorHandler,
		},
	}
	//Инициализируем сервер
	system := ewa.Suffix{
		Index: 2,
		Value: ":system",
	}
	version := ewa.Suffix{
		Index: 3,
		Value: ":version",
	}
	ws, _ := ewa.New("Example", cfg)
	ws.RegisterWeb(new(web.Home), "/")
	ws.RegisterWeb(new(web.Login), "/login")
	ws.RegisterWeb(new(web.Logout), "/logout")
	ws.RegisterWeb(new(section1.Section_1_1), "/section1/1_1")
	ws.RegisterWeb(new(section1.Section_1_2), "/section1/1_2")
	ws.RegisterRest(new(api.User), "", "person", system, version)
	//ws.SetBasicAuth(ba)
	//Cors = nil - DefaultConfig
	ws.SetCors(nil)
	//ws.SetStore(nil)
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
