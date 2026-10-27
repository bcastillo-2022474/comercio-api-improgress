-- name: InsertOutboxRecord :exec
INSERT INTO outbox (id, event_name, payload, created_at)
VALUES (@id, @event_name, @payload, @created_at);
