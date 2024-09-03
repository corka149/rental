package app

import (
	"context"
	"github.com/corka149/rental"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine, ctx context.Context, queries *datastore.Queries, config *rental.Config) {
	router.StaticFS("/static", http.FS(static.Assets))
}
