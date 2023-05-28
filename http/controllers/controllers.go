package controllers

import (
	"sigomid/core/http"
	"sigomid/core/http/utils"
)

func RegisteredControllers() []http.ControllerInterface {

	return []http.ControllerInterface{
		PingController{&utils.DefaultController{}},
		LoginController{&utils.DefaultController{}},
		HomeController{&utils.DefaultController{}},
	}
}
