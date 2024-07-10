-- name: CreateSeminar :one
INSERT INTO seminar (id, created_at, updated_at, name, user_id,expiry_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllSeminars :many
SELECT * FROM seminar WHERE user_id = $1;

-- name: GetAllSeminarsByAPIKey :many
SELECT * FROM seminar WHERE api_key = $1;

-- name: DeleteSeminar :exec
DELETE FROM seminar WHERE id = $1;

-- name: EditSeminarName :one
UPDATE seminar
SET name = $2, updated_at = $3
WHERE id = $1
RETURNING *;


-- name: GetSeminarByName :many
SELECT * FROM seminar WHERE name LIKE $1 and user_id = $2;

-- name: DeleteAfterTwoDays :exec
DELETE FROM seminar
WHERE expiry_date <= NOW();
