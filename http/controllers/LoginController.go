package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct{}

type credentials struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (c LoginController) Register(publicRoutes *gin.RouterGroup, securedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/form", gin.H{})
	})

	publicRoutes.POST("/login", func(c *gin.Context) {

		var userInput credentials
		err := c.Bind(&userInput)
		if err != nil {
			panic(err)
		}
		if userInput.Username != "aaa" || userInput.Password != "bbb" {
			session := sessions.Default(c)
			session.AddFlash("invalid username or password")
			c.Redirect(http.StatusUnauthorized, "/login")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "/securedPing")
		return

	})
}
