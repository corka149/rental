-- +goose Up
-- +goose StatementBegin
CREATE TABLE holidays
(
    id SERIAL PRIMARY KEY,
    beginning DATE NOT NULL,
    ending   DATE NOT NULL,
    title VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS holidays;
-- +goose StatementEnd
