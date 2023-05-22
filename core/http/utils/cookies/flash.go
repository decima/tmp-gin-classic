package cookies

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(Flash{})
}

type level string

const (
	Info    level = "info"
	Warning level = "warning"
	Error   level = "error"
)

type Flash struct {
	Type    level
	Message string
}

func (f Flash) String() string {
	return fmt.Sprintf("%s: %s", f.Type, f.Message)
}

func AddFlash(c *gin.Context, t level, msg string) error {
	session := sessions.Default(c)
	session.AddFlash(Flash{t, msg})
	//session.AddFlash(msg)
	return session.Save()
}

func AddInfoFlash(c *gin.Context, msg string) error {
	return AddFlash(c, Info, msg)
}

func AddWarningFlash(c *gin.Context, msg string) error {
	return AddFlash(c, Warning, msg)
}

func AddErrorFlash(c *gin.Context, msg string) error {
	return AddFlash(c, Error, msg)
}

func Flashes(c *gin.Context) ([]Flash, error) {
	session := sessions.Default(c)
	flashes := session.Flashes()

	var result []Flash
	for _, f := range flashes {
		result = append(result, f.(Flash))
	}
	if err := session.Save(); err != nil {
		return nil, err
	}
	return result, nil
}
