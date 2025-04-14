-- create activities table
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_id VARCHAR(255) NOT NULL UNIQUE,
    actor VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    object_id VARCHAR(255),
    object_type VARCHAR(50),
    target VARCHAR(255),
    raw_content JSONB NOT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- create index for activities table
CREATE INDEX IF NOT EXISTS idx_activities_actor ON activities(actor);
CREATE INDEX IF NOT EXISTS idx_activities_processed ON activities(processed);

