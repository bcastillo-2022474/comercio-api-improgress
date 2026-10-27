// Package shared holds the cross-module application primitives: the unit of
// work and the domain event publisher.
package shared

import (
	"context"
	"encoding/json"
)

// Event is a domain fact. Payloads carry ids only, never personal data.
type Event interface {
	Name() string
	Payload() json.RawMessage
}

// Publisher records events. It never fails: an event is only saved if the unit of work commits.
type Publisher interface {
	Publish(event Event)
}

// UnitOfWork runs fn in one transaction and hands it R, the repositories bound to that transaction.
// Every command goes through it. Events published inside fn are saved to the outbox in the same
// transaction, and only if fn returns nil.
type UnitOfWork[R any] interface {
	Execute(ctx context.Context, fn func(ctx context.Context, r R, pub Publisher) error) error
}

// CollectingPublisher buffers events until the unit of work decides to save them.
type CollectingPublisher struct {
	events []Event
}

func NewCollectingPublisher() *CollectingPublisher {
	return &CollectingPublisher{}
}

func (p *CollectingPublisher) Publish(event Event) {
	p.events = append(p.events, event)
}

func (p *CollectingPublisher) Events() []Event {
	return p.events
}
