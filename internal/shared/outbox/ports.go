package outbox

import "context"

type Repository interface {
	// Save stores records. The unit of work calls it inside the write transaction.
	Save(ctx context.Context, records []Record) error
}
