// Package outbox stores domain events inside the write transaction, so an event exists if and only if its write committed.
// Delivery (the relay and the broker) is not built yet.
package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"joao-api-comercio/internal/shared"
)

type Record struct {
	Id        uuid.UUID
	EventName string
	Payload   json.RawMessage
	CreatedAt time.Time
}

// NewRecord turns an event into an outbox record. The caller supplies the clock.
func NewRecord(event shared.Event, now time.Time) Record {
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}
	return Record{
		Id:        id,
		EventName: event.Name(),
		Payload:   event.Payload(),
		CreatedAt: now,
	}
}
