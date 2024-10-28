-- +goose Up
-- +goose StatementBegin
ALTER TABLE rentals ADD COLUMN street VARCHAR(255);
ALTER TABLE rentals ADD COLUMN city VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rentals DROP COLUMN street;
ALTER TABLE rentals DROP COLUMN city;
-- +goose StatementEnd
