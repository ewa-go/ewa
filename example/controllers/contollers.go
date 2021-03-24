package controllers

import (
	"github.com/egovorukhin/egowebapi"
	"github.com/egovorukhin/egowebapi/example/controllers/api"
	"github.com/egovorukhin/egowebapi/example/controllers/web"
)

func Api() egowebapi.Controllers {
	return egowebapi.Controllers{
		api.NewUser("").Controller,
	}
}

func Web() egowebapi.Controllers {
	return egowebapi.Controllers{
		web.NewIndex("/").Controller,
	}
}
