package controllers

import "github.com/gin-gonic/gin"

type PingController struct{}

func (c PingController) Register(publicRoutes *gin.RouterGroup, securedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/ping", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	securedRoutes.GET("/securedPing", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
}
