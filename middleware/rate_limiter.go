package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(1, 2)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		c.Status(429)
		c.Abort()
	}
}
