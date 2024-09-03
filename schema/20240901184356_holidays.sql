-- +goose Up
-- +goose StatementBegin
CREATE TABLE holidays
(
    id SERIAL PRIMARY KEY,
    "from" DATE NOT NULL,
    "to"   DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS holidays;
-- +goose StatementEnd
