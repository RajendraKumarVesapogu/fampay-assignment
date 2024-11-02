package controllers

import (
	"fampay-assignment/config"
	"fampay-assignment/lib"
	"fampay-assignment/logger"
	"fampay-assignment/services"
	types "fampay-assignment/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetLatestVideos(
	ctx *gin.Context,
	db *pgxpool.Pool,
) (interface{}, error) {
	name := "GetLatestVideos"

	var data types.GetLatestVideosRequest
	data.SortOrder = ctx.Query("sort_order")
	data.PaginationPage, _ = strconv.Atoi(ctx.Query("pagination_page"))
	data.PaginationSize, _ = strconv.Atoi(ctx.Query("pagination_size"))
	data.PublishedAfter = ctx.Query("published_after")
	if data.PublishedAfter == "" {
		data.PublishedAfter = time.Now().Add(-20 * time.Hour).Format(config.DATE_FORMAT)
	}
	err := data.Validate()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"controller": name,
			"data":       data,
			"err":        err,
		}).Error("invalid request")
		return lib.ApiResponse{}, lib.NewExternalError().BadRequest(err.Error())
	}
	res, err := services.GetLatestVideos(db,&data)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"controller": name,
			"err":        err,
		}).Error("error getting latest videos")
		return lib.ApiResponse{}, err
	}
	return res, nil

}
// func GetSegmentAcrossModules(
// 	ctx *gin.Context,
// 	db *pgxpool.Pool,
// ) (interface{}, error) {
// 	name := "GetSegmentAcrossModules"

// 	var data externalTypes.GetSegmentAcrossModulesRequest
// 	ctx.BindJSON(&data)
// 	err := data.Validate()
// 	if err != nil {
// 		logger.Log.WithFields(logrus.Fields{
// 			"controller": name,
// 			"data":       data,
// 			"err":        err,
// 		}).Error("invalid request")
// 		return lib.ApiResponse{}, lib.NewExternalError().BadRequest(err.Error())
// 	}

// 	res, err := services.GetSegmentAcrossModules(db, &data)
// 	if err != nil {
// 		logger.Log.WithFields(logrus.Fields{
// 			"controller": name,
// 			"err":        err,
// 		}).Error("error getting segment across modules")
// 		return lib.ApiResponse{}, err
// 	}

// 	return res, nil
// }
