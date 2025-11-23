DROP INDEX IF EXISTS idx_events_unpublished;
ALTER TABLE events DROP COLUMN IF EXISTS published_at;
