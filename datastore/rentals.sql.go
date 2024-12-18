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
INSERT INTO rentals ("object_id", "beginning", "ending", "description", "street", "city") VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, object_id, beginning, ending, description, street, city
`

type CreateRentalParams struct {
	ObjectID    int32
	Beginning   pgtype.Date
	Ending      pgtype.Date
	Description pgtype.Text
	Street      pgtype.Text
	City        pgtype.Text
}

func (q *Queries) CreateRental(ctx context.Context, arg CreateRentalParams) (Rental, error) {
	row := q.db.QueryRow(ctx, createRental,
		arg.ObjectID,
		arg.Beginning,
		arg.Ending,
		arg.Description,
		arg.Street,
		arg.City,
	)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.Beginning,
		&i.Ending,
		&i.Description,
		&i.Street,
		&i.City,
	)
	return i, err
}

const deleteRental = `-- name: DeleteRental :one
DELETE FROM rentals WHERE id = $1 RETURNING id, object_id, beginning, ending, description, street, city
`

func (q *Queries) DeleteRental(ctx context.Context, id int32) (Rental, error) {
	row := q.db.QueryRow(ctx, deleteRental, id)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.Beginning,
		&i.Ending,
		&i.Description,
		&i.Street,
		&i.City,
	)
	return i, err
}

const deleteRentalsInRange = `-- name: DeleteRentalsInRange :many
DELETE FROM rentals WHERE (beginning BETWEEN $1 AND $2) OR (ending BETWEEN $1 AND $2) RETURNING id, object_id, beginning, ending, description, street, city
`

type DeleteRentalsInRangeParams struct {
	Beginning pgtype.Date
	Ending    pgtype.Date
}

func (q *Queries) DeleteRentalsInRange(ctx context.Context, arg DeleteRentalsInRangeParams) ([]Rental, error) {
	rows, err := q.db.Query(ctx, deleteRentalsInRange, arg.Beginning, arg.Ending)
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
			&i.Beginning,
			&i.Ending,
			&i.Description,
			&i.Street,
			&i.City,
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

const getRentalById = `-- name: GetRentalById :one
SELECT id, object_id, beginning, ending, description, street, city FROM rentals WHERE id = $1
`

func (q *Queries) GetRentalById(ctx context.Context, id int32) (Rental, error) {
	row := q.db.QueryRow(ctx, getRentalById, id)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.Beginning,
		&i.Ending,
		&i.Description,
		&i.Street,
		&i.City,
	)
	return i, err
}

const getRentals = `-- name: GetRentals :many
SELECT id, object_id, beginning, ending, description, street, city FROM rentals ORDER BY "beginning"
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
			&i.Beginning,
			&i.Ending,
			&i.Description,
			&i.Street,
			&i.City,
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

const getRentalsInRangeAllObject = `-- name: GetRentalsInRangeAllObject :many
SELECT id, object_id, beginning, ending, description, street, city FROM rentals WHERE ((beginning BETWEEN $1 AND $2) OR (ending BETWEEN $1 AND $2)) AND id <> $3 ORDER BY beginning
`

type GetRentalsInRangeAllObjectParams struct {
	Beginning pgtype.Date
	Ending    pgtype.Date
	Ignoreid  int32
}

func (q *Queries) GetRentalsInRangeAllObject(ctx context.Context, arg GetRentalsInRangeAllObjectParams) ([]Rental, error) {
	rows, err := q.db.Query(ctx, getRentalsInRangeAllObject, arg.Beginning, arg.Ending, arg.Ignoreid)
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
			&i.Beginning,
			&i.Ending,
			&i.Description,
			&i.Street,
			&i.City,
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

const getRentalsInRangeByObject = `-- name: GetRentalsInRangeByObject :many
SELECT id, object_id, beginning, ending, description, street, city FROM rentals WHERE ((beginning BETWEEN $1 AND $2) OR (ending BETWEEN $1 AND $2)) AND id <> $3 AND object_id = $4 ORDER BY beginning
`

type GetRentalsInRangeByObjectParams struct {
	Beginning pgtype.Date
	Ending    pgtype.Date
	Ignoreid  int32
	Objectid  int32
}

func (q *Queries) GetRentalsInRangeByObject(ctx context.Context, arg GetRentalsInRangeByObjectParams) ([]Rental, error) {
	rows, err := q.db.Query(ctx, getRentalsInRangeByObject,
		arg.Beginning,
		arg.Ending,
		arg.Ignoreid,
		arg.Objectid,
	)
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
			&i.Beginning,
			&i.Ending,
			&i.Description,
			&i.Street,
			&i.City,
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
UPDATE rentals SET "object_id" = $1, "beginning" = $2, "ending" = $3, "description" = $4, "street" = $5, "city" = $6 WHERE id = $7 RETURNING id, object_id, beginning, ending, description, street, city
`

type UpdateRentalParams struct {
	ObjectID    int32
	Beginning   pgtype.Date
	Ending      pgtype.Date
	Description pgtype.Text
	Street      pgtype.Text
	City        pgtype.Text
	ID          int32
}

func (q *Queries) UpdateRental(ctx context.Context, arg UpdateRentalParams) (Rental, error) {
	row := q.db.QueryRow(ctx, updateRental,
		arg.ObjectID,
		arg.Beginning,
		arg.Ending,
		arg.Description,
		arg.Street,
		arg.City,
		arg.ID,
	)
	var i Rental
	err := row.Scan(
		&i.ID,
		&i.ObjectID,
		&i.Beginning,
		&i.Ending,
		&i.Description,
		&i.Street,
		&i.City,
	)
	return i, err
}
