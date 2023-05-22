package controllers

import (
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils"
	"sigomid/core/http/utils/view"
)

type HomeController struct {
	*utils.DefaultController
}

func (h HomeController) Register(publicRoutes, securedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/", landingPage)
}

func landingPage(c *gin.Context) {
	view.Render(c, 200, "home/index", gin.H{})
}
