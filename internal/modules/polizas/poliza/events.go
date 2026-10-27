package poliza

import "encoding/json"

// PolizaCreatedEvent says a declaration exists. Ids only.
type PolizaCreatedEvent struct {
	id int64
}

func NewPolizaCreatedEvent(p Poliza) PolizaCreatedEvent {
	return PolizaCreatedEvent{id: p.Id}
}

func (PolizaCreatedEvent) Name() string { return "poliza.created" }

func (e PolizaCreatedEvent) Payload() json.RawMessage {
	payload, _ := json.Marshal(struct {
		PolizaId int64 `json:"poliza_id"`
	}{e.id})
	return payload
}
