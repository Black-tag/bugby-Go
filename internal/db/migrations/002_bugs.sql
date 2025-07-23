-- +goose Up
CREATE TABLE bugs (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    posted_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (posted_by) REFERENCES users(id)
);

-- +goose Down
DROP TABLES bugs;