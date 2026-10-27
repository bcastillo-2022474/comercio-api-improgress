package outbox

import "context"

type Repository interface {
	// Save stores records. The unit of work calls it inside the write transaction.
	Save(ctx context.Context, records []Record) error
	// Count reports how many records are stored. Tests use it to check what was written.
	Count(ctx context.Context) (int, error)
}
