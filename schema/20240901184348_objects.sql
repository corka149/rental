-- +goose Up
-- +goose StatementBegin
CREATE TABLE objects
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS objects;
-- +goose StatementEnd
