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
SELECT * FROM rentals WHERE (("beginning" <= $1 AND $1 <= "ending") OR ("beginning" <= $2 AND $2 <= "ending")) AND id <> $3 AND object_id = $4::int ORDER BY "beginning"; 

-- name: GetRentalsInRangeAllObject :many
SELECT * FROM rentals WHERE (("beginning" <= $1 AND $1 <= "ending") OR ("beginning" <= $2 AND $2 <= "ending")) AND id <> $3 ORDER BY "beginning"; 
