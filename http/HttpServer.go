package http

import (
	cryptorand "crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type RoutesRegistrar interface {
	Register(publicRoutes *gin.RouterGroup, securedRoutes *gin.RouterGroup)
}
type ServerConfiguration struct {
	HostAndPort string
	Security    ServerSecurityConfiguration
}

type ServerSecurityConfiguration struct {
	Enable                     bool
	RedirectOnUnauthorizedPath string
}

type Server struct {
	Configuration  ServerConfiguration
	sessionStorage sessions.Store
}

func (s *Server) Start(routesRegistrar ...RoutesRegistrar) error {

	buf := make([]byte, 128)
	cryptorand.Read(buf)
	s.sessionStorage = cookie.NewStore(buf)

	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	r.Use(sessions.Sessions("session", s.sessionStorage))
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	r.LoadHTMLGlob("templates/*")
	sessionSecuredRoutes := r.Group("")

	sessionSecuredRoutes.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		if !s.Configuration.Security.Enable {
			c.Next()
			return
		}

		user := session.Get("user")
		if user != nil {
			c.Next()
			return
		}

		c.Redirect(302, s.Configuration.Security.RedirectOnUnauthorizedPath)

	})
	publicRoutes := r.Group("")
	for _, controller := range routesRegistrar {
		controller.Register(publicRoutes, sessionSecuredRoutes)
	}

	return r.Run(s.Configuration.HostAndPort)
}
