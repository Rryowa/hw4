-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    storage_until TIMESTAMPTZ NOT NULL,
    issued BOOLEAN NOT NULL,
    issued_at TIMESTAMPTZ,
    returned BOOLEAN NOT NULL,
    hash VARCHAR(255) NOT NULL
);

CREATE INDEX user_id_storage_desc ON orders (user_id, storage_until DESC);
CREATE INDEX id_sort ON orders (id ASC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX user_id_storage_desc;
DROP INDEX id_sort;
DROP TABLE orders;
-- +goose StatementEnd