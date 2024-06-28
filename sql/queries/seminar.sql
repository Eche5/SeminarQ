-- name: CreateSeminar :one
INSERT INTO seminar (id, created_at, updated_at, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;