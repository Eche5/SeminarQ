-- name: CreateQuestion :one

INSERT INTO question (id, created_at, updated_at, seminar_id, question)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllQuestion :many
SELECT * FROM question WHERE seminar_id = $1;