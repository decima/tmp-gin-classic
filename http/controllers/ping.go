package controllers

import (
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils/controllers"
)

type PingController struct {
	*controllers.DefaultController `route:"/ping"`
	Secured                        gin.HandlerFunc `route:"/secured" method:"GET" group:"ROLE_USER"`
	Unsecured                      gin.HandlerFunc `route:"" method:"GET"`
}

func init() {
	controllers.Register(&PingController{
		DefaultController: &controllers.DefaultController{},
		Secured: func(c *gin.Context) {
			c.JSON(200, "ok")
		},
		Unsecured: func(c *gin.Context) {
			c.JSON(204, nil)
		},
	})
}
