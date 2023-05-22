package security

import (
	"encoding/gob"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(User{})
	gob.Register(gin.H{})
}

type User struct {
	Username string
}

const sessionKey = "user"

func IsLogged(c *gin.Context) bool {
	session := sessions.Default(c)
	return session.Get(sessionKey) != nil
}

func LoggedUser(c *gin.Context) any {
	session := sessions.Default(c)
	u := session.Get(sessionKey)
	if u != nil {
		return u
	}
	return nil
}

func TryLogin(c *gin.Context, username string, password string) error {
	if !IsCsrfValid(c, "login") {
		return errors.New("invalid csrf token")
	}
	if username != "admin" || password != "admin" {
		return errors.New("invalid credentials")
	}

	return login(c, User{username})

}

func login(c *gin.Context, user User) error {
	session := sessions.Default(c)
	session.Set("user", user.Username)
	return session.Save()
}

func Logout(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(sessionKey)
	return session.Save()
}
