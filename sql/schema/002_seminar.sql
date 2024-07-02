-- +goose Up
CREATE TABLE seminar (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    name TEXT NOT NULL,
    api_key VARCHAR(6) UNIQUE NOT NULL DEFAULT lpad(floor(random() * 1000000)::text, 6, '0'),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS seminar;
