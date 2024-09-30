package app

import (
	"context"
	"log"
	"time"

	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/templates"
	"github.com/jackc/pgx/v5/pgtype"
)

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

func rentalConflictsByObject(queries *datastore.Queries, ctx context.Context, excludeId int32, beginning, ending time.Time, objectId int32) templates.ErrorCode {
	queryParam := datastore.GetRentalsInRangeByObjectParams{
		Beginning:   pgtype.Date{Time: beginning, Valid: true},
		Beginning_2: pgtype.Date{Time: ending, Valid: true},
		ID:          excludeId,
		Column4:     objectId,
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
		Beginning:   pgtype.Date{Time: beginning, Valid: true},
		Beginning_2: pgtype.Date{Time: ending, Valid: true},
		ID:          excludeId,
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

func endingConflicts(beginning, ending time.Time) templates.ErrorCode {
	if beginning.After(ending) {
		return templates.ErrEndingBeforeBeginning
	}

	return ""
}
