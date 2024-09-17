package middleware

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthCheck(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("user") == nil {
			c.Redirect(302, "/auth/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
