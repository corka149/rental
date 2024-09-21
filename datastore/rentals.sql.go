// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rentals.sql

package datastore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRental = `-- name: CreateRental :one
INSERT INTO rentals ("object_id", "from", "to", "description") VALUES ($1, $2, $3, $4) RETURNING id, object_id, "from", "to", description
`

type CreateRentalParams struct {
	ObjectID    int32
	From        pgtype.Date
	To          pgtype.Date
	Description pgtype.Text
}

func (q *Queries) CreateRental(ctx context.Context, arg CreateRentalParams) (Rental, error) {
	row := q.db.QueryRow(ctx, createRental,
		arg.ObjectID,
		arg.From,
		arg.To,
		arg.Description,
	)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.From,
		&i.To,
		&i.Description,
	)
	return i, err
}

const getRentalById = `-- name: GetRentalById :one
SELECT id, object_id, "from", "to", description FROM rentals WHERE id = $1
`

func (q *Queries) GetRentalById(ctx context.Context, id int32) (Rental, error) {
	row := q.db.QueryRow(ctx, getRentalById, id)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.From,
		&i.To,
		&i.Description,
	)
	return i, err
}

const getRentals = `-- name: GetRentals :many
SELECT id, object_id, "from", "to", description FROM rentals ORDER BY "from"
`

func (q *Queries) GetRentals(ctx context.Context) ([]Rental, error) {
	rows, err := q.db.Query(ctx, getRentals)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Rental
	for rows.Next() {
		var i Rental
		if err := rows.Scan(
			&i.ID,
			&i.ObjectID,
			&i.From,
			&i.To,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRental = `-- name: UpdateRental :one
UPDATE rentals SET "object_id" = $1, "from" = $2, "to" = $3, "description" = $4 WHERE id = $5 RETURNING id, object_id, "from", "to", description
`

type UpdateRentalParams struct {
	ObjectID    int32
	From        pgtype.Date
	To          pgtype.Date
	Description pgtype.Text
	ID          int32
}

func (q *Queries) UpdateRental(ctx context.Context, arg UpdateRentalParams) (Rental, error) {
	row := q.db.QueryRow(ctx, updateRental,
		arg.ObjectID,
		arg.From,
		arg.To,
		arg.Description,
		arg.ID,
	)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.From,
		&i.To,
		&i.Description,
	)
	return i, err
}
