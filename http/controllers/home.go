package controllers

import (
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils/controllers"
	"sigomid/core/http/utils/view"
)

type HomeController struct {
	*controllers.DefaultController
	Index gin.HandlerFunc `route:"/" method:"GET"`
}

func init() {
	controllers.Register(&HomeController{
		DefaultController: &controllers.DefaultController{},
		Index: func(c *gin.Context) {
			view.Render(c, 200, "home/index", gin.H{})
		},
	})
}
