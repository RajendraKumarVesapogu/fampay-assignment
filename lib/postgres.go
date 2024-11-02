package lib

import (
	"fmt"
	"time"

	"fampay-assignment/connections"
	"fampay-assignment/logger"
	"fampay-assignment/models"
	"fampay-assignment/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


func executePostgresQuery(
	db *pgxpool.Pool,
	queryName string,
	query string,
	queryArgs ...interface{},
) (
	pgx.Rows,
	error,
) {
	
	rows, err := db.Query(connections.GetContext(), query, queryArgs...)
	if err != nil {
		logger.Log.WithFields(
			logger.Fields{
				"query": query,
				"args":  queryArgs,
			},
		).Errorf("Error executing query %s: %v", queryName, err)
	}
	return rows, err
}

type GetLatestYouTubeVideoQueryParams struct {
	PaginationPage int
	PaginationSize int
	PublishedAfter time.Time
	SortOrder	   string
}

type GetLatestYouTubeVideoQueryResult struct {
	Videos []models.Video
	Err    error
}

func GetLatestYouTubeVideoQuery(
	db *pgxpool.Pool,
	params *GetLatestYouTubeVideoQueryParams,
) (response GetLatestYouTubeVideoQueryResult) {
	response.Videos = []models.Video{}

	query := fmt.Sprintf(
		`SELECT
			video_id, title, description, published_at, thumbnail_url, channel_title, channel_id
		FROM
			videos
		WHERE
			published_at > $1
		ORDER BY
			published_at %s
		LIMIT $2 OFFSET $3`,
		params.SortOrder,
	)

	rows, err := executePostgresQuery(
		db,
		"GetLatestYouTubeVideoQuery",
		query,
		params.PublishedAfter,
		params.PaginationSize,
		utils.GetPaginationOffset(params.PaginationPage, params.PaginationSize),
	)
	if err != nil {
		response.Err = err
		return response
	}
	defer rows.Close()

	for rows.Next() {
		var video models.Video
		err := rows.Scan(
			&video.VideoID,
			&video.Title,
			&video.Description,
			&video.PublishedAt,
			&video.ThumbnailURL,
			&video.ChannelTitle,
			&video.ChannelID,
		)
		if err != nil {
			response.Err = err
			return response
		}
		response.Videos = append(response.Videos, video)
	}

	return response
}
