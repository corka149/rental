package app

import (
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
)

func indexHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		templates.Layout(templates.Index()).Render(c.Request.Context(), c.Writer)
	}
}
