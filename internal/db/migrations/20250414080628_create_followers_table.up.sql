-- create followers table
CREATE TABLE IF NOT EXISTS followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_actor_id VARCHAR(255) NOT NULL,
    follower_inbox VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, follower_actor_id)
);

-- create index for followers table
CREATE INDEX IF NOT EXISTS idx_followers_user_id ON followers(user_id);
