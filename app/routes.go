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

	// ==================== OBJECT ====================
	router.GET("/objects", indexObjects(queries))
	router.GET("/objects/new", newObjectForm(queries))
	router.POST("/objects/new", createObject(queries))
	router.GET("/objects/:id", updateObjectForm(queries))
	router.POST("/objects/:id", updateObject(queries))
	router.POST("/objects/:id/delete", deleteObject(queries))

	// ==================== HOLIDAY ====================
	router.GET("/holidays", indexHolidays(queries))
	router.GET("/holidays/new", newHolidayForm(queries))
	router.POST("/holidays/new", createHoliday(queries))
	router.GET("/holidays/:id", updateHolidayForm(queries))
	router.POST("/holidays/:id", updateHoliday(queries))
	router.POST("/holidays/:id/delete", deleteHoliday(queries))
}
