package main

import (
	"fmt"
	ewa "github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
	"os"
	"strings"
)

func main() {

	//WEB
	cfg := ewa.Config{
		Port:    3003,
		Timeout: ewa.NewTimeout(30, 30, 30),
		Views: &ewa.Views{
			Root:   "www",
			Engine: ".html",
		},
		Static: "www",
	}
	//BasicAuth
	/*users := map[string]string{
		"user": "Qq123456",
	}
	authorizer := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}
	ba := ewa.NewBasicAuth(users, authorizer, nil)*/
	//Инициализируем сервер
	ws, _ := ewa.New("Example", cfg)
	ws.RegisterWeb(new(web.Index), "/")
	ws.RegisterRest(new(api.User), "")
	//ws.SetBasicAuth(ba)
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
