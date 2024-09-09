package middleware

import (
	"cmp"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
)

func NewI18n() gin.HandlerFunc {

	return func(c *gin.Context) {
		lang := cmp.Or(c.Query("lang"), c.GetHeader("Accept-Language"), "en")

		ctx, err := ctxi18n.WithLocale(c.Request.Context(), lang)

		if err != nil {
			log.Println(err)
		} else {
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}
