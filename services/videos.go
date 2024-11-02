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
	videosResult := lib.GetLatestYouTubeVideoQuery(
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


// func GetSegments(
// 	db *pgxpool.Pool,
// 	params *externalTypes.GetSegmentsRequest,
// ) (
// 	segments externalTypes.GetSegmentsResponse,
// 	err error,
// ) {
	

// 	moduleConfig, ok := findModuleConfig(params.Module, &orgConfig)
// 	if !ok || !moduleConfig.Enabled {
// 		errorMessage := "module not found or disabled"
// 		logger.Log.WithFields(logger.Fields{
// 			"params": params,
// 		}).Error(errorMessage)
// 		return segments, lib.NewExternalError().BadRequest(errorMessage)
// 	}

// 	sortColumn, ok := getSortColumn(moduleConfig, params.SortField, params.SortMetricName)
// 	if !ok {
// 		errorMessage := "invalid sort params"
// 		logger.Log.WithFields(logger.Fields{
// 			"params": params,
// 			"user":   user,
// 		}).Error(sortColumn)
// 		return segments, lib.NewExternalError().BadRequest(errorMessage)
// 	}

// 	postDate, _ := utils.ParseInputDtmString(params.PostDate)
// 	result := getSegmentDerivedMetrics(
// 		user,
// 		db,
// 		&GetSegmentDerivedMetricsParams{
// 			Module: moduleConfig,

// 			PreDataColumnSuffix:  params.PrePeriod,
// 			PostDate:             postDate,
// 			PostDataColumnSuffix: params.PostPeriod,

// 			PaginationPage: params.PaginationPage,
// 			PaginationSize: params.PaginationSize,

// 			SortColumn:   &sortColumn,
// 			SortOrder:    &params.SortOrder,
// 			SortAbsolute: &params.SortAbsolute,
// 			Props:        params.Props,
// 			Vals:         params.Vals,
// 			PropCount:    params.PropCount,
// 			Starred:      params.Starred,
// 			Filter:       params.Filter,
// 		},
// 	)
// 	segments = result.Segments
// 	err = result.Err

// 	if segments == nil {
// 		segments = externalTypes.GetSegmentsResponse{}
// 	}
// 	return segments, err
// }