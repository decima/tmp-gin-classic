package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sigomid/core/http/utils"
	"sigomid/core/http/utils/cookies"
	"sigomid/core/http/utils/security"
	"sigomid/core/http/utils/view"
	"sigomid/http/dataObjects"
)

type LoginController struct {
	*utils.DefaultController
}

func (controller LoginController) Register(publicRoutes, securedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/login", controller.Login)
	publicRoutes.POST("/login", controller.Login)
	publicRoutes.GET("/logout", controller.Logout)
}

func (controller LoginController) Login(c *gin.Context) {
	var userInput dataObjects.LoginForm
	c.Bind(&userInput)

	if c.Request.Method == http.MethodPost {
		if err := security.TryLogin(c, userInput.Username, userInput.Password); err != nil {
			log.Println(err)
			_ = cookies.AddErrorFlash(c, "Invalid username or password")
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.Redirect(http.StatusFound, controller.ServerConfiguration.Security.RedirectOnLogin)
		return

	}

	view.Render(c, http.StatusOK, "login/form", gin.H{
		"username": userInput.Username,
	})

	return

}

func (controller LoginController) Logout(c *gin.Context) {
	_ = security.Logout(c)
	c.Redirect(http.StatusFound, "/login")
}
