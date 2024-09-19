-- name: GetHolidays :many
SELECT * FROM holidays ORDER BY beginning;

-- name: GetHolidayById :one
SELECT * FROM holidays WHERE id = $1;

-- name: CreateHoliday :one
INSERT INTO holidays (beginning, ending, title) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateHoliday :one
UPDATE holidays SET beginning = $1, ending = $2, title = $3 WHERE id = $4 RETURNING *;

-- name: DeleteHoliday :one
DELETE FROM holidays WHERE id = $1 RETURNING *;