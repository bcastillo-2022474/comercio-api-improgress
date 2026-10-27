// Package fakes holds in-memory stand-ins for shared infrastructure.
package fakes

import (
	"context"
	"time"

	"joao-api-comercio/internal/shared"
	"joao-api-comercio/internal/shared/outbox"
)

// UnitOfWork runs the callback against fixed repositories. When the callback succeeds it saves
// the published events to the outbox, as the real one does in its transaction. A failing callback saves nothing.
type UnitOfWork[R any] struct {
	repos  R
	outbox outbox.Repository
	now    func() time.Time
}

func NewUnitOfWork[R any](repos R, outbox outbox.Repository, now func() time.Time) *UnitOfWork[R] {
	return &UnitOfWork[R]{repos: repos, outbox: outbox, now: now}
}

func (u *UnitOfWork[R]) Execute(ctx context.Context, fn func(ctx context.Context, r R, pub shared.Publisher) error) error {
	pub := shared.NewCollectingPublisher()
	if err := fn(ctx, u.repos, pub); err != nil {
		return err
	}

	records := make([]outbox.Record, 0, len(pub.Events()))
	for _, event := range pub.Events() {
		records = append(records, outbox.NewRecord(event, u.now()))
	}
	return u.outbox.Save(ctx, records)
}
