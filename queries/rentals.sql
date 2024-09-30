-- name: GetRentals :many
SELECT * FROM rentals ORDER BY "from";

-- name: GetRentalById :one
SELECT * FROM rentals WHERE id = $1;

-- name: CreateRental :one
INSERT INTO rentals ("object_id", "from", "to", "description") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateRental :one
UPDATE rentals SET "object_id" = $1, "from" = $2, "to" = $3, "description" = $4 WHERE id = $5 RETURNING *;

-- name: DeleteRental :one
DELETE FROM rentals WHERE id = $1 RETURNING *;

-- name: GetRentalsInRangeByObject :many
SELECT * FROM rentals WHERE (("from" <= $1 AND $1 <= "to") OR ("from" <= $2 AND $2 <= "to")) AND id <> $3 AND object_id = $4::int ORDER BY "from"; 

-- name: GetRentalsInRangeAllObject :many
SELECT * FROM rentals WHERE (("from" <= $1 AND $1 <= "to") OR ("from" <= $2 AND $2 <= "to")) AND id <> $3 ORDER BY "from"; 
