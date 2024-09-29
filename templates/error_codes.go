package templates

import "fmt"

type ErrorCode string

const (
	// General
	ErrUnableToGetData ErrorCode = "unable_to_get_data"

	// Holiday
	ErrHolidayConfictsWithAnother ErrorCode = "holiday_conflicts_with_another"
)

func (e ErrorCode) String() string {
	return fmt.Sprintf("errors.%s", string(e))
}
