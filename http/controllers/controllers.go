package controllers

import (
	"sigomid/core/http"
	"sigomid/core/http/utils"
)

func RegisteredControllers() []http.RoutesRegistrar {

	return []http.RoutesRegistrar{
		PingController{&utils.DefaultController{}},
		LoginController{&utils.DefaultController{}},
		HomeController{&utils.DefaultController{}},
	}
}
