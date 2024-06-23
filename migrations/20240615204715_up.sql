-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    storage_until TIMESTAMPTZ NOT NULL,
    issued BOOLEAN NOT NULL,
    issued_at TIMESTAMPTZ,
    returned BOOLEAN NOT NULL,
    order_price FLOAT NOT NULL,
    weight FLOAT NOT NULL,
    package_type VARCHAR(255) NOT NULL,
    package_price FLOAT NOT NULL,
    hash VARCHAR(255) NOT NULL
);

CREATE INDEX user_id_storage_asc ON orders (user_id, storage_until ASC);
CREATE INDEX id_asc ON orders (id ASC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX user_id_storage_asc;
DROP INDEX id_asc;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd