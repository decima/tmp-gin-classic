package utils

import "sigomid/core/http"

type DefaultController struct {
	ServerConfiguration http.ServerConfiguration
}

func (d DefaultController) SetConfiguration(s http.ServerConfiguration) {
	d.ServerConfiguration = s
}
