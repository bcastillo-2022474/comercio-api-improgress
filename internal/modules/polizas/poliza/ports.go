package poliza

import (
	"context"

	"joao-api-comercio/internal/shared/types"
)

type Repository interface {
	// Insert stores the declaration and all its items, and sets p.Id. A correlativo already in use
	// fails with *CorrelativoTakenError and writes nothing. Several rows are written, so it belongs
	// inside the caller's unit of work, which makes them atomic.
	Insert(ctx context.Context, p *Poliza) error
	// GetById returns the declaration with its items in the order they were stored. An unknown id is
	// absent, not an error.
	GetById(ctx context.Context, id int64) (types.Maybe[Poliza], error)
}
