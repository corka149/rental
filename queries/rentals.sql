-- name: GetRentals :many
SELECT * FROM rentals ORDER BY "beginning";

-- name: GetRentalById :one
SELECT * FROM rentals WHERE id = $1;

-- name: CreateRental :one
INSERT INTO rentals ("object_id", "beginning", "ending", "description") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateRental :one
UPDATE rentals SET "object_id" = $1, "beginning" = $2, "ending" = $3, "description" = $4 WHERE id = $5 RETURNING *;

-- name: DeleteRental :one
DELETE FROM rentals WHERE id = $1 RETURNING *;

-- name: GetRentalsInRangeByObject :many
SELECT * FROM rentals WHERE ((beginning BETWEEN sqlc.arg(beginning) AND sqlc.arg(ending)) OR (ending BETWEEN sqlc.arg(beginning) AND sqlc.arg(ending))) AND id <> sqlc.arg(ignoreId) AND object_id = sqlc.arg(objectId) ORDER BY beginning; 

-- name: GetRentalsInRangeAllObject :many
SELECT * FROM rentals WHERE ((beginning BETWEEN sqlc.arg(beginning) AND sqlc.arg(ending)) OR (ending BETWEEN sqlc.arg(beginning) AND sqlc.arg(ending))) AND id <> sqlc.arg(ignoreId) ORDER BY beginning;
