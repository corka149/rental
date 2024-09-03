-- name: GetHolidays :many
SELECT * FROM holidays;

-- name: GetHolidayById :one
SELECT * FROM holidays WHERE id = $1;

-- name: CreateHoliday :one
INSERT INTO holidays ("from", "to") VALUES ($1, $2) RETURNING *;

-- name: UpdateHoliday :one
UPDATE holidays SET "from" = $1, "to" = $2 WHERE id = $3 RETURNING *;

-- name: DeleteHoliday :one
DELETE FROM holidays WHERE id = $1 RETURNING *;