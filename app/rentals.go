package app

import (
	"log"
	"slices"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
)

func indexRentals(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		rentals, err := queries.GetRentals(c.Request.Context())

		if err != nil {
			log.Printf("Error getting rentals: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		objectIds := make([]int32, 10)

		for _, r := range rentals {
			if slices.Contains(objectIds, r.ObjectID) {
				continue
			}

			objectIds = append(objectIds, r.ObjectID)
		}

		objects, err := queries.GetObjectByIds(c.Request.Context(), objectIds)

		if err != nil {
			log.Printf("Error getting objects: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		objectsById := make(map[int32]datastore.Object, len(objects))

		for _, o := range objects {
			objectsById[o.ID] = o
		}

		templates.Layout(user.Name, templates.RentalIndex(rentals, objectsById)).Render(c.Request.Context(), c.Writer)
	}
}
