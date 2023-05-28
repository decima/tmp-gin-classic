package security

import (
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils/configuration"
	"sigomid/core/http/utils/controllers"
)

func RoleMiddlewareResolver(configuration configuration.ServerConfiguration, requestedGroup string) func(c *gin.Context) {

	return func(c *gin.Context) {
		if requestedGroup != controllers.PUBLIC_GROUP {
			if !IsLogged(c) && configuration.Security.Enable {
				c.Redirect(302, configuration.Security.RedirectOnUnauthorizedPath)
				return
			}
			if !LoggedUser(c).HasRole(requestedGroup) {
				c.JSON(403, "forbidden")
				return
			}
		}
		c.Next()
		return

	}

}
