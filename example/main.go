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

	//BasicAuth
	authorizer := func(user string, pass string) bool {
		if user == "user" && pass == "Qq123456" {
			return true
		}
		return false
	}
	//WEB
	cfg := ewa.Config{
		Port:    3005,
		Timeout: ewa.NewTimeout(30, 30, 30),
		Views: &ewa.Views{
			Root:   "www",
			Engine: ".html",
		},
		Static:    "www",
		BasicAuth: ewa.NewBasicAuth(authorizer, nil),
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
	ws.RegisterWeb(new(web.Index), "/")
	ws.RegisterRest(new(api.User), "", "person", system, version)
	//ws.SetBasicAuth(ba)
	//Cors = nil - DefaultConfig
	ws.SetCors(nil)
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
