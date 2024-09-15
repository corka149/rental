package app

import (
	"context"
	"net/http"

	"github.com/corka149/rental"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/static"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, ctx context.Context, queries *datastore.Queries, config *rental.Config) {
	router.StaticFS("/static", http.FS(static.Assets))

	// ==================== HOME ====================
	router.GET("/", indexHome())

	// ==================== AUTH ====================
	auth := router.Group("/auth")
	auth.POST("/register", register(queries))
	auth.GET("/login", loginForm())
	auth.POST("/login", login(queries))
	auth.GET("/logout", logout())
}
