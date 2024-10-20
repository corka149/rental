package jobs

import (
	"context"
	"log"
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/jackc/pgx/v5/pgtype"
)

func CleanUp(ctx context.Context, queries *datastore.Queries) {
	go doCleanUp(ctx, queries)
}

func doCleanUp(ctx context.Context, queries *datastore.Queries) {
	t := time.Tick(24 * time.Hour)

	for next := range t {
		RemoveOldEntries(next, ctx, queries)
	}
}

func RemoveOldEntries(anchor time.Time, ctx context.Context, queries *datastore.Queries) {
	deleted := 0
	minusOneYear := anchor.AddDate(-1, 0, 0)
	minueOneDay := anchor.AddDate(0, 0, -1)

	holidays, err := queries.DeleteHolidaysInRange(ctx, datastore.DeleteHolidaysInRangeParams{
		Beginning: pgtype.Date{Time: minusOneYear, Valid: true},
		Ending:    pgtype.Date{Time: minueOneDay, Valid: true},
	})

	deleted += len(holidays)

	if err != nil {
		log.Printf("Error deleting holidays: %v\n", err)
	}

	_, err = queries.DeleteRentalsInRange(ctx, datastore.DeleteRentalsInRangeParams{
		Beginning: pgtype.Date{Time: minusOneYear, Valid: true},
		Ending:    pgtype.Date{Time: minueOneDay, Valid: true},
	})

	if err != nil {
		log.Printf("Error deleting rentals: %v\n", err)
	}

	deleted += len(holidays)

	log.Printf("Deleted %d holidays and rentals\n", deleted)
}
