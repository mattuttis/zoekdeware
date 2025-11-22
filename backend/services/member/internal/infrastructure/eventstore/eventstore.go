package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	ErrStreamNotFound    = errors.New("event stream not found")
	ErrConcurrencyConflict = errors.New("concurrency conflict")
)

type StoredEvent struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	AggregateID string          `json:"aggregate_id"`
	Data        json.RawMessage `json:"data"`
	Metadata    json.RawMessage `json:"metadata"`
	Version     int             `json:"version"`
	Timestamp   time.Time       `json:"timestamp"`
}

type EventStore interface {
	Append(ctx context.Context, aggregateID string, expectedVersion int, events []StoredEvent) error
	Load(ctx context.Context, aggregateID string) ([]StoredEvent, error)
	LoadFrom(ctx context.Context, aggregateID string, fromVersion int) ([]StoredEvent, error)
}

type InMemoryEventStore struct {
	mu      sync.RWMutex
	streams map[string][]StoredEvent
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		streams: make(map[string][]StoredEvent),
	}
}

func (s *InMemoryEventStore) Append(ctx context.Context, aggregateID string, expectedVersion int, events []StoredEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, exists := s.streams[aggregateID]
	if !exists {
		stream = make([]StoredEvent, 0)
	}

	currentVersion := len(stream)
	if expectedVersion != -1 && currentVersion != expectedVersion {
		return ErrConcurrencyConflict
	}

	for i, event := range events {
		event.Version = currentVersion + i + 1
		event.Timestamp = time.Now()
		stream = append(stream, event)
	}

	s.streams[aggregateID] = stream
	return nil
}

func (s *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]StoredEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stream, exists := s.streams[aggregateID]
	if !exists {
		return nil, ErrStreamNotFound
	}

	result := make([]StoredEvent, len(stream))
	copy(result, stream)
	return result, nil
}

func (s *InMemoryEventStore) LoadFrom(ctx context.Context, aggregateID string, fromVersion int) ([]StoredEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stream, exists := s.streams[aggregateID]
	if !exists {
		return nil, ErrStreamNotFound
	}

	if fromVersion >= len(stream) {
		return []StoredEvent{}, nil
	}

	result := make([]StoredEvent, len(stream)-fromVersion)
	copy(result, stream[fromVersion:])
	return result, nil
}
