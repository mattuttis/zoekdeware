package messaging

import (
	"context"
	"encoding/json"
	"time"
)

type Message struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Metadata    MessageMetadata `json:"metadata"`
	PublishedAt time.Time       `json:"published_at"`
}

type MessageMetadata struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	Source        string `json:"source,omitempty"`
}

type Publisher interface {
	Publish(ctx context.Context, topic string, message Message) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(ctx context.Context, message Message) error

type MessageBroker interface {
	Publisher
	Subscriber
}
