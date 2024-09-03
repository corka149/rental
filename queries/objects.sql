-- name: GetObjects :many
SELECT * FROM objects;

-- name: GetObjectById :one
SELECT * FROM objects WHERE id = $1;

-- name: CreateObject :one
INSERT INTO objects (name) VALUES ($1) RETURNING *;

-- name: UpdateObject :one
UPDATE objects SET name = $1 WHERE id = $2 RETURNING *;

-- name: DeleteObject :one
DELETE FROM objects WHERE id = $1 RETURNING *;