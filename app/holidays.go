package app

import (
	"fmt"
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

		target := fmt.Sprintf("%d", holiday.ID)

		templates.Layout(user.Name, templates.HolidayForm(holiday, target)).Render(c.Request.Context(), c.Writer)
	}
}

func updateHoliday(queries *datastore.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
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
