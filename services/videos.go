package services

import (
	"time"

	"fampay-assignment/config"
	"fampay-assignment/lib"
	"fampay-assignment/logger"
	types "fampay-assignment/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetLatestVideos(
	db *pgxpool.Pool,
	params *types.GetLatestVideosRequest,
) (
	response types.GetLatestVideosResponse,
	err error,
) {

	err = params.Validate()
	if err != nil {
		logger.Log.WithFields(logger.Fields{
			"params": params,
		}).Error(err)
		return response, lib.NewExternalError().BadRequest(err.Error())
	}
	publishedAfter, _ := time.Parse(config.DATE_FORMAT,params.PublishedAfter)
	videosResult := lib.GetLatestYouTubeVideos(
		db,
		&lib.GetLatestYouTubeVideoQueryParams{
			SortOrder: params.SortOrder,
			PaginationSize: params.PaginationSize,
			PaginationPage: params.PaginationPage,
			PublishedAfter: publishedAfter,
		},
	)

	if videosResult.Err != nil{
		logger.Log.WithFields(
			logger.Fields{
				"params": params,
			},	
		).Error(err)
	}
	response.Videos = videosResult.Videos
	if response.Videos == nil {
		response = types.GetLatestVideosResponse{}
	}
	return response, err
}