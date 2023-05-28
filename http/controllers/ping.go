package controllers

import (
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils"
)

type PingController struct {
	*utils.DefaultController `route:"/ping"`
}

func (controller PingController) Register(publicRoutes, securedRoutes *gin.RouterGroup) {
	publicRoutes.GET("", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	securedRoutes.GET("/secured", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
}
