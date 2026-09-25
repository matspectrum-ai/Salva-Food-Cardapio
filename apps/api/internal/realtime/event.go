package realtime

import (
	"context"
	"encoding/json"
	"time"
)

const (
	EventOrderCreated = "order.created"
	EventOrderUpdated = "order.updated"
)

type Event struct {
	ID            string
	TenantID      string
	AggregateType string
	AggregateID   string
	Type          string
	Payload       []byte
	OccurredAt    time.Time
}

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID            string          `json:"id"`
		TenantID      string          `json:"tenant_id"`
		AggregateType string          `json:"aggregate_type"`
		AggregateID   string          `json:"aggregate_id"`
		Type          string          `json:"type"`
		Payload       json.RawMessage `json:"payload"`
		OccurredAt    time.Time       `json:"occurred_at"`
	}{
		ID: e.ID, TenantID: e.TenantID, AggregateType: e.AggregateType,
		AggregateID: e.AggregateID, Type: e.Type, Payload: e.Payload,
		OccurredAt: e.OccurredAt,
	})
}

type Bus interface {
	Publish(context.Context, Event) error
	Subscribe(context.Context, string) (<-chan Event, error)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var wire struct {
		ID            string          `json:"id"`
		TenantID      string          `json:"tenant_id"`
		AggregateType string          `json:"aggregate_type"`
		AggregateID   string          `json:"aggregate_id"`
		Type          string          `json:"type"`
		Payload       json.RawMessage `json:"payload"`
		OccurredAt    time.Time       `json:"occurred_at"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	e.ID, e.TenantID, e.AggregateType, e.AggregateID = wire.ID, wire.TenantID, wire.AggregateType, wire.AggregateID
	e.Type, e.Payload, e.OccurredAt = wire.Type, append([]byte(nil), wire.Payload...), wire.OccurredAt
	return nil
}
