-- Event store table for member service
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(36) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    metadata JSONB,
    version INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(aggregate_id, version)
);

-- Index for loading aggregate event streams
CREATE INDEX idx_events_aggregate ON events(aggregate_id, version);

-- Index for event type queries (for projections)
CREATE INDEX idx_events_type ON events(event_type, created_at);
