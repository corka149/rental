package app

import (
	"context"
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
			From: now,
			To:   now,
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
		fromStr := c.PostForm("from")
		toStr := c.PostForm("to")
		description := c.PostForm("description")
		objectIDStr := c.PostForm("object")
		objectID, err := strconv.Atoi(objectIDStr)

		if err != nil {
			log.Printf("Error parsing object id: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		from, err := time.Parse("2006-01-02", fromStr)

		if err != nil {
			log.Printf("Error parsing from date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		to, err := time.Parse("2006-01-02", toStr)

		if err != nil {
			log.Printf("Error parsing to date: %v", err)
			c.Redirect(302, "/rentals")
			c.Abort()
			return
		}

		errCodes := make([]templates.ErrorCode, 0)

		if errCode := rentalConflictsByObject(queries, c.Request.Context(), 0, from, to, int32(objectID)); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := holidayConflicts(queries, c.Request.Context(), 0, from, to); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if len(errCodes) > 0 {
			rental := datastore.Rental{
				From:        pgtype.Date{Time: from},
				To:          pgtype.Date{Time: to},
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
			From:        pgtype.Date{Time: from, Valid: true},
			To:          pgtype.Date{Time: to, Valid: true},
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

func rentalConflictsByObject(queries *datastore.Queries, ctx context.Context, excludeId int32, beginning, ending time.Time, objectId int32) templates.ErrorCode {
	queryParam := datastore.GetRentalsInRangeByObjectParams{
		From:    pgtype.Date{Time: beginning, Valid: true},
		From_2:  pgtype.Date{Time: ending, Valid: true},
		ID:      excludeId,
		Column4: objectId,
	}

	rentals, err := queries.GetRentalsInRangeByObject(ctx, queryParam)

	if err != nil {
		log.Printf("Error getting rentals: %v", err)
		return templates.ErrUnableToGetData
	}

	if len(rentals) > 0 {
		return templates.ErrConflictsWithRental
	}

	return ""
}

func rentalConflicts(queries *datastore.Queries, ctx context.Context, excludeId int32, beginning, ending time.Time) templates.ErrorCode {
	queryParam := datastore.GetRentalsInRangeAllObjectParams{
		From:   pgtype.Date{Time: beginning, Valid: true},
		From_2: pgtype.Date{Time: ending, Valid: true},
		ID:     excludeId,
	}

	rentals, err := queries.GetRentalsInRangeAllObject(ctx, queryParam)

	if err != nil {
		log.Printf("Error getting rentals: %v", err)
		return templates.ErrUnableToGetData
	}

	if len(rentals) > 0 {
		return templates.ErrConflictsWithRental
	}

	return ""
}
