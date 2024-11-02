package models

import (
	"time"
)

type Video struct {
	VideoID      string    `db:"video_id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	PublishedAt  time.Time `db:"published_at"`
	ThumbnailURL string    `db:"thumbnail_url"`
	ChannelTitle string    `db:"channel_title"`
	ChannelID    string    `db:"channel_id"`
}
