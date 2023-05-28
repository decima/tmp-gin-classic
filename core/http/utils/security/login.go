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

var roleHierarchy = map[string][]string{
	"ROLE_ADMIN": {"ROLE_USER"},
	"ROLE_USER":  {},
}

type User struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func (u User) String() string {

	return u.Username
}

func (u User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
		for _, rr := range roleHierarchy[r] {
			if rr == role {
				return true
			}
		}

	}
	return false
}

const sessionKey = "user"

func IsLogged(c *gin.Context) bool {
	return LoggedUser(c).Username != ""
}

func LoggedUser(c *gin.Context) User {
	session := sessions.Default(c)
	u := session.Get(sessionKey)
	if u != nil {
		return u.(User)
	}
	return User{}
}

func TryLogin(c *gin.Context, username string, password string) error {
	if !IsCsrfValid(c, "login") {
		return errors.New("invalid csrf token")
	}
	if username != "admin" || password != "admin" {
		return errors.New("invalid credentials")
	}

	return login(c, User{Username: username, Roles: []string{"ROLE_ADMIN"}})
}

func login(c *gin.Context, user User) error {
	session := sessions.Default(c)
	session.Set("user", user)

	return session.Save()
}

func Logout(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(sessionKey)

	return session.Save()
}
