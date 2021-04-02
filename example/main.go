package main

import (
	"fmt"
	"github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
	"os"
	"strings"
)

func main() {

	//API
	cfgApi := egowebapi.Config{
		Port:    3003,
		Timeout: egowebapi.NewTimeout(30, 30, 30),
	}
	rest, _ := egowebapi.New("api", cfgApi)
	rest.SetRest(new(api.User), "")
	rest.Start()

	//WEB
	cfgWeb := egowebapi.Config{
		Port:    3000,
		Timeout: egowebapi.NewTimeout(30, 30, 30),
		Views: &egowebapi.Views{
			Root: "www",
			Ext:  ".html",
		},
		Static: "www",
	}
	http, _ := egowebapi.New("http", cfgWeb)
	http.SetWeb(new(web.Index), "")
	http.Start()

	for {
		var input string
		_, err := fmt.Fscan(os.Stdin, &input)
		if err != nil {
			os.Exit(1)
		}
		switch strings.ToLower(input) {
		case "exit":
			//fmt.Println(api.Stop())
			fmt.Println(http.Stop())
			os.Exit(0)
		}
	}
}
