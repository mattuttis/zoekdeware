-- Add publishing tracking to events table for outbox pattern
ALTER TABLE events ADD COLUMN published_at TIMESTAMPTZ;

-- Index for efficient polling of unpublished events
CREATE INDEX idx_events_unpublished ON events(id) WHERE published_at IS NULL;
