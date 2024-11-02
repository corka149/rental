-- +goose Up
-- +goose StatementBegin
ALTER TABLE rentals ADD COLUMN house_number VARCHAR(255);
ALTER TABLE rentals ADD COLUMN postal_code VARCHAR(255);
ALTER TABLE rentals ADD COLUMN with_delivery BOOLEAN;
ALTER TABLE rentals ADD COLUMN with_setup BOOLEAN;
ALTER TABLE rentals ADD COLUMN beginning_time TIME;
ALTER TABLE rentals ADD COLUMN ending_time TIME;
ALTER TABLE rentals ADD COLUMN price_per_object DECIMAL(10, 2);
ALTER TABLE rentals ADD COLUMN total_price DECIMAL(10, 2);
ALTER TABLE rentals ADD COLUMN notes TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rentals DROP COLUMN house_number;
ALTER TABLE rentals DROP COLUMN postal_code;
ALTER TABLE rentals DROP COLUMN with_delivery;
ALTER TABLE rentals DROP COLUMN with_setup;
ALTER TABLE rentals DROP COLUMN beginning_time;
ALTER TABLE rentals DROP COLUMN ending_time;
ALTER TABLE rentals DROP COLUMN price_per_object;
ALTER TABLE rentals DROP COLUMN total_price;
ALTER TABLE rentals DROP COLUMN notes;
-- +goose StatementEnd
