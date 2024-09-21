-- name: GetRentals :many
SELECT * FROM rentals ORDER BY "from";

-- name: GetRentalById :one
SELECT * FROM rentals WHERE id = $1;

-- name: CreateRental :one
INSERT INTO rentals ("object_id", "from", "to", "description") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateRental :one
UPDATE rentals SET "object_id" = $1, "from" = $2, "to" = $3, "description" = $4 WHERE id = $5 RETURNING *;
