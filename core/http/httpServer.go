package http

import (
	cryptorand "crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"sigomid/core/http/utils/configuration"
	"sigomid/core/http/utils/controllers"
	"sigomid/core/http/utils/cookies"
	"sigomid/core/http/utils/security"

	"text/template"

	_ "sigomid/http/controllers"
)

type Server struct {
	Configuration  *configuration.ServerConfiguration
	sessionStorage sessions.Store
}

func (s *Server) Start() error {

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

	for _, controller := range controllers.GetComputed(*s.Configuration) {

		for _, action := range controller.Actions {
			roleResolverRouter := r.Group("", security.RoleMiddlewareResolver(*s.Configuration, action.Group))
			if action.Method == controllers.ANY_METHOD {
				roleResolverRouter.Any(action.Route, action.Do)
				continue
			}
			roleResolverRouter.Handle(action.Method, action.Route, action.Do)
		}

	}

	return r.Run(s.Configuration.HostAndPort)
}
