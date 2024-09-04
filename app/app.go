package app

import (
	"context"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
)

func indexHome(ctx context.Context, queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {

		templates.Layout(templates.Index()).Render(ctx, c.Writer)
	}
}
