package templates

import "fmt"

type ErrorCode string

const (
	// General
	ErrUnableToGetData ErrorCode = "unable_to_get_data"

	// Holiday
	ErrConflictsWithHoliday ErrorCode = "conflicts_with_holiday"

	// Rental
	ErrConflictsWithRental ErrorCode = "conflicts_with_rental"
)

func (e ErrorCode) String() string {
	return fmt.Sprintf("errors.%s", string(e))
}
