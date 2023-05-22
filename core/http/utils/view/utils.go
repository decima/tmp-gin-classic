package view

import (
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, statusCode int, template string, params gin.H) {
	globals := gin.H{
		"Context": c,
	}

	for k, v := range params {
		globals[k] = v
	}

	c.HTML(statusCode, template, globals)

}
