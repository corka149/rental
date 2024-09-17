// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package datastore

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Holiday struct {
	ID   int32
	From pgtype.Date
	To   pgtype.Date
}

type Object struct {
	ID   int32
	Name string
}

type Rental struct {
	ID          int32
	ObjectID    int32
	From        pgtype.Date
	To          pgtype.Date
	Description pgtype.Text
}

type User struct {
	ID       int32
	Email    string
	Password string
	Name     string
}
