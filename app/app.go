package app

import (
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
)

func indexHome(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)

		templates.Layout(user.Name, templates.Index()).Render(c.Request.Context(), c.Writer)
	}
}
