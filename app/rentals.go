package app

import (
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

func newRentalForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		now := pgtype.Date{
			Time: time.Now(),
		}

		rental := datastore.Rental{
			Beginning: now,
			Ending:    now,
		}

		objects, err := queries.GetObjects(c.Request.Context())

		if err != nil {
			log.Printf("Error getting objects: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.RentalForm(rental, "new", objects)).Render(c.Request.Context(), c.Writer)
	}
}

func createRental(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		beginningStr := c.PostForm("beginning")
		toStr := c.PostForm("ending")
		description := c.PostForm("description")
		objectIDStr := c.PostForm("object")
		objectID, err := strconv.Atoi(objectIDStr)

		if err != nil {
			log.Printf("Error parsing object id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		beginning, err := time.Parse("2006-01-02", beginningStr)

		if err != nil {
			log.Printf("Error parsing beginning date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		ending, err := time.Parse("2006-01-02", toStr)

		if err != nil {
			log.Printf("Error parsing to date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		errCodes := make([]templates.ErrorCode, 0)

		if errCode := rentalConflictsByObject(queries, c.Request.Context(), 0, beginning, ending, int32(objectID)); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := holidayConflicts(queries, c.Request.Context(), 0, beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := endingConflicts(beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if len(errCodes) > 0 {
			rental := datastore.Rental{
				Beginning:   pgtype.Date{Time: beginning},
				Ending:      pgtype.Date{Time: ending},
				ObjectID:    int32(objectID),
				Description: pgtype.Text{String: description},
			}

			objects, err := queries.GetObjects(c.Request.Context())

			if err != nil {
				log.Printf("Error getting objects: %v", err)
				c.Redirect(302, "/rentals")
				c.Abort()
				return
			}

			templates.Layout(user.Name, templates.RentalForm(rental, "new", objects, errCodes...)).Render(c.Request.Context(), c.Writer)
			return
		}

		rental := datastore.CreateRentalParams{
			Beginning:   pgtype.Date{Time: beginning, Valid: true},
			Ending:      pgtype.Date{Time: ending, Valid: true},
			Description: pgtype.Text{String: description, Valid: true},
			ObjectID:    int32(objectID),
		}

		_, err = queries.CreateRental(c.Request.Context(), rental)

		if err != nil {
			log.Printf("Error creating rental: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		c.Redirect(302, "/rentals")
	}
}

func updateRentalForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing rental id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		rental, err := queries.GetRentalById(c.Request.Context(), int32(id))

		if err != nil {
			log.Printf("Error getting rental: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		objects, err := queries.GetObjects(c.Request.Context())

		if err != nil {
			log.Printf("Error getting objects: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.RentalForm(rental, idStr, objects)).Render(c.Request.Context(), c.Writer)
	}
}

func updateRental(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing rental id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		beginningStr := c.PostForm("beginning")
		toStr := c.PostForm("ending")
		description := c.PostForm("description")
		objectIDStr := c.PostForm("object")
		objectID, err := strconv.Atoi(objectIDStr)

		if err != nil {
			log.Printf("Error parsing object id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		beginning, err := time.Parse("2006-01-02", beginningStr)

		if err != nil {
			log.Printf("Error parsing beginning date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		ending, err := time.Parse("2006-01-02", toStr)

		if err != nil {
			log.Printf("Error parsing to date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		errCodes := make([]templates.ErrorCode, 0)

		if errCode := endingConflicts(beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := rentalConflictsByObject(queries, c.Request.Context(), int32(id), beginning, ending, int32(objectID)); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := holidayConflicts(queries, c.Request.Context(), int32(id), beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if len(errCodes) > 0 {
			rental := datastore.Rental{
				ID:          int32(id),
				Beginning:   pgtype.Date{Time: beginning},
				Ending:      pgtype.Date{Time: ending},
				ObjectID:    int32(objectID),
				Description: pgtype.Text{String: description},
			}

			objects, err := queries.GetObjects(c.Request.Context())

			if err != nil {
				log.Printf("Error getting objects: %v", err)
				c.Redirect(302, "/rentals")
				c.Abort()
				return
			}

			templates.Layout(user.Name, templates.RentalForm(rental, idStr, objects, errCodes...)).Render(c.Request.Context(), c.Writer)
			return
		}

		rental := datastore.UpdateRentalParams{
			ID:          int32(id),
			Beginning:   pgtype.Date{Time: beginning, Valid: true},
			Ending:      pgtype.Date{Time: ending, Valid: true},
			Description: pgtype.Text{String: description, Valid: true},
			ObjectID:    int32(objectID),
		}

		_, err = queries.UpdateRental(c.Request.Context(), rental)

		if err != nil {
			log.Printf("Error updating rental: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		c.Redirect(302, "/rentals")
	}
}

func deleteRental(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing rental id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		_, err = queries.DeleteRental(c.Request.Context(), int32(id))

		if err != nil {
			log.Printf("Error deleting rental: %v", err)
		}

		c.Redirect(302, "/rentals")
	}
}
