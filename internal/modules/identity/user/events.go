package user

import (
	"encoding/json"

	"joao-api-comercio/internal/shared/types"
)

// UserRegisteredEvent says a new account exists. Ids only: the email stays out of the event.
type UserRegisteredEvent struct {
	id types.UserId
}

func NewUserRegisteredEvent(u User) UserRegisteredEvent {
	return UserRegisteredEvent{id: u.Id}
}

func (UserRegisteredEvent) Name() string { return "user.registered" }

func (e UserRegisteredEvent) Payload() json.RawMessage {
	payload, _ := json.Marshal(struct {
		UserId string `json:"user_id"`
	}{e.id.String()})
	return payload
}
