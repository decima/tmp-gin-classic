package http

import (
	cryptorand "crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils/cookies"
	"sigomid/core/http/utils/security"
	"text/template"
)

type RoutesRegistrar interface {
	SetConfiguration(ServerConfiguration)
	Register(publicRoutes, securedRoutes *gin.RouterGroup)
}
type ServerConfiguration struct {
	HostAndPort string
	Security    ServerSecurityConfiguration
}

type ServerSecurityConfiguration struct {
	Enable                      bool
	RedirectOnUnauthorizedPath  string
	RedirectOnLogin             string
	UsePreviousIfDefinedOnLogin bool
}

type Server struct {
	Configuration  *ServerConfiguration
	sessionStorage sessions.Store
}

func (s *Server) Start(routesRegistrar ...RoutesRegistrar) error {

	buf := make([]byte, 128)
	cryptorand.Read(buf)
	//s.sessionStorage = cookie.NewStore(buf)

	var err error

	s.sessionStorage, err = redis.NewStore(
		10,
		"tcp",
		"localhost:16379",
		"eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		buf,
	)
	s.sessionStorage.Options(sessions.Options{
		MaxAge: 86400, // 1 day
		Path:   "/",
	})

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("public", false)))

	r.Use(sessions.Sessions("session", s.sessionStorage))
	r.Use(func(c *gin.Context) {
		c.Next()
	})
	r.SetFuncMap(template.FuncMap{
		"csrf":       security.Csrf,
		"flashes":    cookies.Flashes,
		"isLogged":   security.IsLogged,
		"loggedUser": security.LoggedUser,
	})
	r.LoadHTMLGlob("templates/*")

	sessionSecuredRoutes := r.Group("")
	sessionSecuredRoutes.Use(s.secureArea)

	publicRoutes := r.Group("")
	for _, controller := range routesRegistrar {
		controller.SetConfiguration(*s.Configuration)
		controller.Register(publicRoutes, sessionSecuredRoutes)
	}

	return r.Run(s.Configuration.HostAndPort)
}

func (s *Server) secureArea(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if s.Configuration.Security.Enable && user == nil {
		c.Redirect(302, s.Configuration.Security.RedirectOnUnauthorizedPath)
		return
	}

	c.Next()

}
