package controllers

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
)

func Api() egowebapi.Controllers {
	return egowebapi.Controllers{
		api.NewUser(""),
	}
}

func Web() egowebapi.Controllers {
	path := "/"
	return egowebapi.Controllers{
		web.NewIndex(path),
	}
}