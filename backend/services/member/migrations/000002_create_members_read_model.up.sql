-- Members read model (denormalized view for queries)
CREATE TABLE IF NOT EXISTS members (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(50),
    bio TEXT,
    birth_date DATE,
    gender VARCHAR(20),
    interests TEXT[], -- PostgreSQL array for interests
    photos TEXT[],    -- PostgreSQL array for photo URLs
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for email lookups
CREATE INDEX idx_members_email ON members(email);

-- Index for status filtering
CREATE INDEX idx_members_status ON members(status);

-- Index for discovery queries (active members)
CREATE INDEX idx_members_discovery ON members(status, gender, birth_date)
    WHERE status = 'active';
