-- +goose Up
CREATE TABLE outbox (
    id         TEXT PRIMARY KEY,
    event_name TEXT NOT NULL,
    payload    TEXT NOT NULL,
    created_at TEXT NOT NULL
);

-- +goose Down
DROP TABLE outbox;
