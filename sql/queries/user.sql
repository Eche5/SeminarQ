
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, full_name,email, password)
VALUES ($1, $2, $3, $4,$5,$6)
RETURNING *;


-- name: GetUserEmail :one
SELECT * FROM users WHERE email= $1;