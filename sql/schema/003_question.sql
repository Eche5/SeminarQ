-- +goose Up
CREATE TABLE question (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seminar_id UUID NOT NULL REFERENCES seminar(id) ON DELETE CASCADE,
    question TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS question;
