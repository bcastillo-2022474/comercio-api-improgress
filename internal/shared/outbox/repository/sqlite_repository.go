// Package repository holds the SQL adapter for the outbox port.
package repository

import (
	"context"
	"fmt"
	"time"

	"joao-api-comercio/internal/shared/outbox"
	"joao-api-comercio/internal/shared/outbox/repository/db"
)

type SQLiteRepository struct {
	q *db.Queries
}

func NewSQLiteRepository(q db.DBTX) *SQLiteRepository {
	return &SQLiteRepository{q: db.New(q)}
}

func (r *SQLiteRepository) Count(ctx context.Context) (int, error) {
	n, err := r.q.CountOutboxRecords(ctx)
	if err != nil {
		return 0, fmt.Errorf("counting outbox records: %w", err)
	}
	return int(n), nil
}

func (r *SQLiteRepository) Save(ctx context.Context, records []outbox.Record) error {
	for _, record := range records {
		err := r.q.InsertOutboxRecord(ctx, db.InsertOutboxRecordParams{
			ID:        record.Id.String(),
			EventName: record.EventName,
			Payload:   string(record.Payload),
			CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
		if err != nil {
			return fmt.Errorf("inserting outbox record: %w", err)
		}
	}
	return nil
}
