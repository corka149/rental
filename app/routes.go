package app

import (
	"context"
	"net/http"

	"github.com/corka149/rental"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/middleware"
	"github.com/corka149/rental/static"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, ctx context.Context, queries *datastore.Queries, config *rental.Config) {
	router.StaticFS("/static", http.FS(static.Assets))

	// ==================== HOME ====================
	router.GET("/", indexHome(queries))

	// ==================== AUTH ====================
	ba := gin.BasicAuth(gin.Accounts{
		config.AdminUsername: config.AdminUserPassword,
	})
	rLimiter := middleware.RateLimiter()

	auth := router.Group("/auth")
	auth.POST("/register", rLimiter, ba, register(queries))
	auth.GET("/login", loginForm(queries))
	auth.POST("/login", rLimiter, login(queries))
	auth.GET("/logout", logout())
}
