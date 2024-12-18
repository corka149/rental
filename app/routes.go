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

	basicAuth := gin.BasicAuth(gin.Accounts{
		config.AdminUsername: config.AdminUserPassword,
	})

	rLimiter := middleware.RateLimiter()

	// ==================== HOME ====================
	router.GET("/", indexHome(queries))

	// ==================== AUTH ====================

	auth := router.Group("/auth")

	auth.POST("/register", rLimiter, basicAuth, register(queries))
	auth.GET("/login", loginForm(queries))
	auth.POST("/login", rLimiter, login(queries))
	auth.GET("/logout", logout())

	// ==================== OBJECT ====================
	objects := router.Group("/objects", middleware.AuthCheck())

	objects.GET("", indexObjects(queries))
	objects.GET("/new", newObjectForm(queries))
	objects.POST("/new", createObject(queries))
	objects.GET("/:id", updateObjectForm(queries))
	objects.POST("/:id", updateObject(queries))
	objects.POST("/:id/delete", deleteObject(queries))

	// ==================== HOLIDAY ====================
	holidays := router.Group("/holidays", middleware.AuthCheck())

	holidays.GET("", indexHolidays(queries))
	holidays.GET("/new", newHolidayForm(queries))
	holidays.POST("/new", createHoliday(queries))
	holidays.GET("/:id", updateHolidayForm(queries))
	holidays.POST("/:id", updateHoliday(queries))
	holidays.POST("/:id/delete", deleteHoliday(queries))

	// ==================== RENTAL ====================
	rentals := router.Group("/rentals", middleware.AuthCheck())

	rentals.GET("", indexRentals(queries))
	rentals.GET("/new", newRentalForm(queries))
	rentals.POST("/new", createRental(queries))
	rentals.GET("/:id", updateRentalForm(queries))
	rentals.POST("/:id", updateRental(queries))
	rentals.POST("/:id/delete", deleteRental(queries))

	// ==================== CALENDAR ====================
	router.GET("/calendar", indexCalendar(queries))
	router.GET("/calendar/search", searchCalendar(queries))

	// ==================== ADMIN ====================
	admin := router.Group("/admin")

	admin.GET("/cleanup", rLimiter, basicAuth, ForceCleanUp(queries))
}
