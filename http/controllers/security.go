package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sigomid/core/http/utils/controllers"
	"sigomid/core/http/utils/cookies"
	"sigomid/core/http/utils/security"
	"sigomid/core/http/utils/view"
	"sigomid/http/dataObjects"
)

type LoginController struct {
	*controllers.DefaultController
	LoginGet  gin.HandlerFunc `route:"/login" method:"GET"`
	LoginPost gin.HandlerFunc `route:"/login" method:"POST"`
	Logout    gin.HandlerFunc `route:"/logout" method:"GET"`
}

func init() {

	controller := &LoginController{
		DefaultController: &controllers.DefaultController{},
	}
	controller.LoginGet = controller.loginAction
	controller.LoginPost = controller.loginAction
	controller.Logout = controller.logoutAction
	controllers.Register(controller)
}

func (controller LoginController) loginAction(c *gin.Context) {
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

func (controller LoginController) logoutAction(c *gin.Context) {
	_ = security.Logout(c)
	c.Redirect(http.StatusFound, "/login")
}
