// Package fakes holds in-memory implementations of the poliza ports for unit tests.
package fakes

import (
	"context"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/shared/types"
)

// Repository enforces the same unique correlativo as the real store.
type Repository struct {
	next        int64
	byId        map[int64]poliza.Poliza
	correlativo map[int64]bool
}

func NewRepository() *Repository {
	return &Repository{byId: map[int64]poliza.Poliza{}, correlativo: map[int64]bool{}}
}

func (r *Repository) Insert(_ context.Context, p *poliza.Poliza) error {
	if r.correlativo[p.Correlativo] {
		return &poliza.CorrelativoTakenError{Correlativo: p.Correlativo}
	}
	r.next++
	p.Id = r.next
	stored := *p
	stored.Items = append([]poliza.Item(nil), p.Items...)
	r.byId[p.Id] = stored
	r.correlativo[p.Correlativo] = true
	return nil
}

func (r *Repository) GetById(_ context.Context, id int64) (types.Maybe[poliza.Poliza], error) {
	p, ok := r.byId[id]
	if !ok {
		return types.None[poliza.Poliza](), nil
	}
	return types.Some(p), nil
}
