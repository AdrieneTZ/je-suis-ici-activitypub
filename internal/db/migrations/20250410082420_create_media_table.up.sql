-- create media table
CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    checkin_id UUID NOT NULL REFERENCES checkins(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INT NOT NULL,
    width INT,
    height INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- create index for media table
CREATE INDEX IF NOT EXISTS idx_media_checkin_id ON media(checkin_id);