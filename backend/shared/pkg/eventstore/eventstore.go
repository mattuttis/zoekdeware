package eventstore

import (
	"context"
	"encoding/json"
	"time"
)

type Event struct {
	ID          string          `json:"id"`
	AggregateID string          `json:"aggregate_id"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	Metadata    Metadata        `json:"metadata"`
	Version     int             `json:"version"`
	Timestamp   time.Time       `json:"timestamp"`
}

type Metadata struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	CausationID   string `json:"causation_id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
}

type EventStore interface {
	Append(ctx context.Context, aggregateID string, expectedVersion int, events []Event) error
	Load(ctx context.Context, aggregateID string) ([]Event, error)
	LoadFrom(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler) error
}
