package app

import (
	"log"

	"github.com/corka149/rental/datastore"
	"github.com/gin-gonic/gin"
)

func IndexObjects(quieres *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := quieres.GetObjects(c.Request.Context())

		if err != nil {
			log.Printf("Error getting objects: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

	}
}
