// Package fakes holds in-memory implementations of the outbox port for unit tests.
package fakes

import (
	"context"
	"sync"

	"joao-api-comercio/internal/shared/outbox"
)

// Repository remembers saved records so a test can assert on them.
type Repository struct {
	mu      sync.Mutex
	records []outbox.Record
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Save(_ context.Context, records []outbox.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, records...)
	return nil
}

func (r *Repository) Records() []outbox.Record {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]outbox.Record(nil), r.records...)
}
