package main

import (
	cryptorand "crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"sigomid/core/http"
	"sigomid/http/controllers"
)

var store sessions.Store

func init() {
	buf := make([]byte, 128)
	cryptorand.Read(buf)
	store = cookie.NewStore(buf)
}

func main() {
	httpServer := http.Server{Configuration: &http.ServerConfiguration{
		HostAndPort: ":8000",
		Security: http.ServerSecurityConfiguration{
			Enable:                      true,
			RedirectOnUnauthorizedPath:  "/login",
			RedirectOnLogin:             "/",
			UsePreviousIfDefinedOnLogin: false}}}

	panic(httpServer.Start(controllers.RegisteredControllers()...))
}
