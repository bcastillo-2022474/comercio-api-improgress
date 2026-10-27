package types

import "github.com/google/uuid"

// IDs are uuid-backed. UUIDv7 is time-ordered, so primary-key inserts stay sequential in the index.
type UserId uuid.UUID

// newUUID returns a v7 UUID, falling back to v4 if the entropy source fails.
func newUUID() uuid.UUID {
	if u, err := uuid.NewV7(); err == nil {
		return u
	}
	return uuid.New()
}

func (id UserId) String() string { return uuid.UUID(id).String() }

func NewUserId() UserId { return UserId(newUUID()) }

func ParseUserId(s string) (UserId, error) {
	u, err := uuid.Parse(s)
	return UserId(u), err
}

func (id UserId) MarshalText() ([]byte, error)  { return uuid.UUID(id).MarshalText() }
func (id *UserId) UnmarshalText(b []byte) error { return (*uuid.UUID)(id).UnmarshalText(b) }
