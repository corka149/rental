-- +goose Up
-- +goose StatementBegin
CREATE TABLE rentals
(
    id SERIAL PRIMARY KEY,
    "object_id"   INTEGER NOT NULL,
    "beginning"        DATE    NOT NULL,
    "ending"          DATE    NOT NULL,
    "description" TEXT,
    FOREIGN KEY (object_id) REFERENCES objects (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rentals;
-- +goose StatementEnd
