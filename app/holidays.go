package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func indexHolidays(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		holidays, err := queries.GetHolidays(c.Request.Context())

		if err != nil {
			log.Printf("Error getting holidays: %v", err)
			c.Redirect(302, "/")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.HolidayIndex(holidays)).Render(c.Request.Context(), c.Writer)
	}
}

func newHolidayForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		now := pgtype.Date{
			Time: time.Now(),
		}
		holiday := datastore.Holiday{
			Beginning: now,
			Ending:    now,
		}
		templates.Layout(user.Name, templates.HolidayForm(holiday, "new")).Render(c.Request.Context(), c.Writer)
	}
}

func createHoliday(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		beginningStr := c.PostForm("beginning")
		endingStr := c.PostForm("ending")
		title := c.PostForm("title")

		beginning, err := time.Parse("2006-01-02", beginningStr)

		if err != nil {
			log.Printf("Error parsing from date: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		ending, err := time.Parse("2006-01-02", endingStr)

		if err != nil {
			log.Printf("Error parsing to date: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		if ending.Before(beginning) {
			log.Printf("Error: to date is before from date")
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		errCodes := make([]templates.ErrorCode, 0)

		if errCode := rentalConflicts(queries, c.Request.Context(), 0, beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := holidayConflicts(queries, c.Request.Context(), 0, beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if len(errCodes) > 0 {
			log.Printf("Error: holiday conflicts with existing holiday")

			holiday := datastore.Holiday{
				Title:     title,
				Beginning: pgtype.Date{Time: beginning, Valid: true},
				Ending:    pgtype.Date{Time: ending, Valid: true},
			}

			templates.Layout(user.Name, templates.HolidayForm(holiday, "new", errCodes...)).Render(c.Request.Context(), c.Writer)
			return
		}

		holiday := datastore.CreateHolidayParams{
			Title:     title,
			Beginning: pgtype.Date{Time: beginning, Valid: true},
			Ending:    pgtype.Date{Time: ending, Valid: true},
		}

		if _, err := queries.CreateHoliday(c.Request.Context(), holiday); err != nil {
			log.Printf("Error creating holiday: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		c.Redirect(302, "/holidays")
	}
}

func updateHolidayForm(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing holiday id: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		holiday, err := queries.GetHolidayById(c.Request.Context(), int32(id))

		if err != nil {
			log.Printf("Error getting holiday: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		templates.Layout(user.Name, templates.HolidayForm(holiday, idStr)).Render(c.Request.Context(), c.Writer)
	}
}

func updateHoliday(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUserFromSession(c, queries)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing holiday id: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		fromStr := c.PostForm("beginning")
		toStr := c.PostForm("ending")
		title := c.PostForm("title")

		beginning, err := time.Parse("2006-01-02", fromStr)

		if err != nil {
			log.Printf("Error parsing from date: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		ending, err := time.Parse("2006-01-02", toStr)

		if err != nil {
			log.Printf("Error parsing to date: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		if ending.Before(beginning) {
			log.Printf("Error: to date is before from date")
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		errCodes := make([]templates.ErrorCode, 0)

		if errCode := rentalConflicts(queries, c.Request.Context(), 0, beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if errCode := holidayConflicts(queries, c.Request.Context(), 0, beginning, ending); errCode != "" {
			errCodes = append(errCodes, errCode)
		}

		if len(errCodes) > 0 {
			log.Printf("Error: holiday conflicts with existing holiday")

			holiday := datastore.Holiday{
				Title:     title,
				Beginning: pgtype.Date{Time: beginning, Valid: true},
				Ending:    pgtype.Date{Time: ending, Valid: true},
			}

			templates.Layout(user.Name, templates.HolidayForm(holiday, idStr, errCodes...)).Render(c.Request.Context(), c.Writer)
			return
		}

		holiday := datastore.UpdateHolidayParams{
			ID:        int32(id),
			Title:     title,
			Beginning: pgtype.Date{Time: beginning, Valid: true},
			Ending:    pgtype.Date{Time: ending, Valid: true},
		}

		if _, err := queries.UpdateHoliday(c.Request.Context(), holiday); err != nil {
			log.Printf("Error updating holiday: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		c.Redirect(302, "/holidays")
	}
}

func deleteHoliday(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Printf("Error parsing holiday id: %v", err)
			c.Redirect(302, "/holidays")
			c.Abort()
			return
		}

		if _, err := queries.DeleteHoliday(c.Request.Context(), int32(id)); err != nil {
			log.Printf("Error deleting holiday: %v", err)
		}

		c.Redirect(302, "/holidays")
	}
}

func holidayConflicts(queries *datastore.Queries, ctx context.Context, excludeId int32, beginning, ending time.Time) templates.ErrorCode {
	queryParams := datastore.GetHolidaysInRangeParams{
		Beginning: pgtype.Date{Time: beginning, Valid: true},
		Ending:    pgtype.Date{Time: ending, Valid: true},
		Ignoreid:  excludeId,
	}

	holidays, err := queries.GetHolidaysInRange(ctx, queryParams)

	if err != nil {
		log.Printf("Error getting holidays in range: %v", err)
		return templates.ErrUnableToGetData
	}

	if len(holidays) > 0 {
		return templates.ErrConflictsWithHoliday
	}

	return ""
}
