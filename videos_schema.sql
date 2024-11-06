CREATE TABLE videos (
    video_id VARCHAR(50) NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ NOT NULL,
    thumbnail_url TEXT,
    channel_title VARCHAR(255),
    channel_id VARCHAR(50) NOT NULL
);

-- Indexes
CREATE INDEX idx_videos_video_id ON videos(video_id);
CREATE INDEX idx_videos_published_at ON videos(published_at);
