package main

import (
	"fmt"
	"github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers"
	"os"
	"strings"
)

func main() {
	/*cfgApi := egowebapi.Config{
		Port:    3000,
		Timeout: egowebapi.NewTimeout(30,30,30),
	}
	api := egowebapi.New("api", cfgApi).SetControllers(controllers.Api())
	api.Start()*/

	cfgWeb := egowebapi.Config{
		Port:    3003,
		Timeout: egowebapi.NewTimeout(30, 30, 30),
	}
	web := egowebapi.New("web", cfgWeb).SetControllers(controllers.Web())
	web.Start()

	for {
		var input string
		_, err := fmt.Fscan(os.Stdin, &input)
		if err != nil {
			os.Exit(1)
		}
		switch strings.ToLower(input) {
		case "exit":
			//fmt.Println(api.Stop())
			fmt.Println(web.Stop())
			break
		}
	}
}
