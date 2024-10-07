package app

import (
	"cmp"
	"log"
	"strconv"
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarEntry struct {
	OccursOn time.Time `json:"occurs_on"`
}

func indexCalendar(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := ctxi18n.Locale(c.Request.Context())

		lang := cmp.Or(locale.Code().String(), "en")

		objects, err := queries.GetObjects(c.Request.Context())

		if err != nil {
			log.Printf("Error fetching objects: %v", err)
			objects = []datastore.Object{}
		}

		templates.CalendarIndex(lang, objects).Render(c.Request.Context(), c.Writer)
	}
}

func searchCalendar(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		objectIdStr := c.Query("object")
		objectId, err := strconv.Atoi(objectIdStr)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid object ID"})
			return
		}

		monthStr := c.Query("month")
		month, err := strconv.Atoi(monthStr)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid month"})
			return
		}

		yearStr := c.Query("year")
		year, err := strconv.Atoi(yearStr)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid year"})
			return
		}

		beginning := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		ending := beginning.AddDate(0, 1, 0)

		params := datastore.GetRentalsInRangeByObjectParams{
			Objectid:  int32(objectId),
			Beginning: pgtype.Date{Time: beginning, Valid: true},
			Ending:    pgtype.Date{Time: ending, Valid: true},
			Ignoreid:  0,
		}

		rentals, err := queries.GetRentalsInRangeByObject(c.Request.Context(), params)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		entriesData := make([]CalendarEntry, 0)

		for _, entry := range rentals {
			occursOn := entry.Beginning.Time

			entriesData = append(entriesData, CalendarEntry{OccursOn: occursOn})
			occursOn = occursOn.AddDate(0, 0, 1)

			for occursOn.Before(entry.Ending.Time) {
				entriesData = append(entriesData, CalendarEntry{OccursOn: occursOn})
				occursOn = occursOn.AddDate(0, 0, 1)
			}

			if entry.Beginning.Time != entry.Ending.Time {
				entriesData = append(entriesData, CalendarEntry{OccursOn: entry.Ending.Time})
			}
		}

		holidayParams := datastore.GetHolidaysInRangeParams{
			Beginning: pgtype.Date{Time: beginning, Valid: true},
			Ending:    pgtype.Date{Time: ending, Valid: true},
			Ignoreid:  0,
		}

		holidays, err := queries.GetHolidaysInRange(c.Request.Context(), holidayParams)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		for _, entry := range holidays {
			occursOn := entry.Beginning.Time

			entriesData = append(entriesData, CalendarEntry{OccursOn: occursOn})

			occursOn = occursOn.AddDate(0, 0, 1)

			for occursOn.Before(entry.Ending.Time) {
				entriesData = append(entriesData, CalendarEntry{OccursOn: occursOn})
				occursOn = occursOn.AddDate(0, 0, 1)
			}

			if entry.Beginning.Time != entry.Ending.Time {
				entriesData = append(entriesData, CalendarEntry{OccursOn: entry.Ending.Time})
			}
		}

		c.JSON(200, entriesData)
	}
}
