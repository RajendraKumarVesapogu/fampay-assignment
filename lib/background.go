package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"fampay-assignment/config"
	"fampay-assignment/connections"
	"fampay-assignment/logger"
	"fampay-assignment/models"
)

const (
	dbOperationTimeout = 30 * time.Second
	maxRetries         = 3
	httpTimeout        = 30 * time.Second
)

var (
	ApiKeys = NewAPIKeys([]string{
		config.YoutubeApiKey1,
		config.YoutubeApiKey2,
		config.YoutubeApiKey3})
)

type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				High struct {
					URL string `json:"url"`
				} `json:"default"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
}

type APIKeys struct {
	keys    []string
	currKey int
}

func (k *APIKeys) nextKey() (string, error) {
	if len(k.keys) == 0 {
		return "", errors.New("no API keys available")
	}
	k.currKey = (k.currKey + 1) % len(k.keys)
	return k.keys[k.currKey], nil
}

func (k *APIKeys) currentKey() string {
	return k.keys[k.currKey]
}

func AddKey(newKey string) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Println("Recovered in AddKey:", r)
		}
	}()

	for _, existingKey := range ApiKeys.keys {
		if existingKey == newKey {
			return true, nil 
		}
	}
	
	if newKey == "" {
		return false,errors.New("new API key is empty")
	}

	ApiKeys.keys = append(ApiKeys.keys, newKey)
	return true, nil 
}

func NewAPIKeys(keys []string) *APIKeys {
	return &APIKeys{keys: keys, currKey: 0}
}

func executeQuery(
	db *pgxpool.Pool,
	query string,
	queryArgs ...interface{},
) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbOperationTimeout)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := db.Exec(ctx, query, queryArgs...)
		if err == nil {
			return nil
		}

		lastErr = err
		if ctx.Err() != nil {
			return fmt.Errorf("context error during database operation: %v", ctx.Err())
		}

		backoffDuration := time.Duration(attempt+1) * 500 * time.Millisecond
		time.Sleep(backoffDuration)
	}

	return fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}

func fetchAndStoreVideos(searchQuery string, db *pgxpool.Pool, publishedAfter time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	for {
		apiKey := ApiKeys.currentKey()
		publishedAfterStr := publishedAfter.UTC().Format("2006-01-02T15:04:05Z")
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&part=snippet&type=video&order=date&q=%s&publishedAfter=%s",
			apiKey, searchQuery, publishedAfterStr)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request to YouTube API: %v", err)
		}
		defer resp.Body.Close()

		var ytResponse YouTubeResponse
		if err := json.NewDecoder(resp.Body).Decode(&ytResponse); err != nil {
			return fmt.Errorf("error decoding YouTube API response: %v", err)
		}

		if ytResponse.Error.Code == 403 && ytResponse.Error.Errors[0].Reason == "quotaExceeded" {
			logger.Log.Printf("Quota exceeded for key. Switching to next key...")
			_, err := ApiKeys.nextKey()
			if err != nil {
				logger.Log.Println("All API keys are exhausted.")
				return err
			}
			continue
		}

		for _, item := range ytResponse.Items {
			video := models.Video{
				VideoID:      item.ID.VideoID,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				PublishedAt:  item.Snippet.PublishedAt,
				ThumbnailURL: item.Snippet.Thumbnails.High.URL,
				ChannelTitle: item.Snippet.ChannelTitle,
				ChannelID:    item.Snippet.ChannelID,
			}

			queryTemplate := `
                INSERT INTO videos (
                    video_id, title, description, published_at, 
                    thumbnail_url, channel_title, channel_id
                )
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (video_id) DO NOTHING`

			err = executeQuery(db, queryTemplate,
				video.VideoID, video.Title, video.Description,
				video.PublishedAt, video.ThumbnailURL,
				video.ChannelTitle, video.ChannelID)

			if err != nil {
				logger.Log.WithFields(logger.Fields{
					"video_id": video.VideoID,
					"error":    err,
				}).Error("Failed to insert video")
				continue
			}

			logger.Log.WithFields(logger.Fields{
				"video_id": video.VideoID,
			}).Info("Successfully inserted video")
		}
		break
	}
	return nil
}

func StartFetchingVideos(ctx context.Context) {
	db, ok := connections.GetPostgresDb()
	if !ok {
		logger.Log.Fatal("Error connecting to database")
		return
	}

	searchQuery := config.YOUTUBE_SEARCH_QUERY
	publishedAfter := time.Now().Add(-100 * time.Minute)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopping video fetch service...")
			return
		case <-ticker.C:
			if err := fetchAndStoreVideos(searchQuery, db, publishedAfter); err != nil {
				logger.Log.WithError(err).Error("Error in fetchAndStoreVideos")
				continue
			}
			publishedAfter = time.Now().Add(10 * time.Second)
		}
	}
}

func init() {
	ctx := context.Background()
	go func() {
		time.Sleep(10 * time.Second)
		StartFetchingVideos(ctx)
	}()
}

