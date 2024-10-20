package app

import (
	"cmp"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
)

func indexHome(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)

		locale := ctxi18n.Locale(c.Request.Context())

		lang := cmp.Or(locale.Code().String(), "en")

		objects, err := queries.GetObjects(c.Request.Context())

		if err != nil {
			objects = []datastore.Object{}
		}

		templates.Layout(user.Name, templates.Index(lang, objects)).Render(c.Request.Context(), c.Writer)
	}
}
