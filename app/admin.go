package app

import (
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/jobs"
	"github.com/gin-gonic/gin"
)

func ForceCleanUp(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		jobs.RemoveOldEntries(now, c.Request.Context(), queries)
	}
}
