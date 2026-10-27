package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"joao-api-comercio/internal/shared"
	"joao-api-comercio/internal/shared/outbox"
	outboxrepo "joao-api-comercio/internal/shared/outbox/repository"
)

// UnitOfWork runs each command in one SQL transaction. The repositories are bound to that
// transaction, and the published events are written to the outbox before commit.
type UnitOfWork[R any] struct {
	db   *sql.DB
	bind func(tx *sql.Tx) R
	now  func() time.Time
}

// NewUnitOfWork takes bind, which builds the repositories R on top of a transaction.
func NewUnitOfWork[R any](db *sql.DB, bind func(tx *sql.Tx) R, now func() time.Time) *UnitOfWork[R] {
	return &UnitOfWork[R]{db: db, bind: bind, now: now}
}

func (u *UnitOfWork[R]) Execute(ctx context.Context, fn func(ctx context.Context, r R, pub shared.Publisher) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	pub := shared.NewCollectingPublisher()
	if err := fn(ctx, u.bind(tx), pub); err != nil {
		return err
	}

	records := make([]outbox.Record, 0, len(pub.Events()))
	for _, event := range pub.Events() {
		records = append(records, outbox.NewRecord(event, u.now()))
	}
	if err := outboxrepo.NewSQLiteRepository(tx).Save(ctx, records); err != nil {
		return fmt.Errorf("saving outbox: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
