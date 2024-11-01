package lib

import (
	"net/http"

	"fampay-assignment/connections"
	"fampay-assignment/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Controller func(
	ctx *gin.Context,
	db *pgxpool.Pool,
) (interface{}, error)

type ApiResponse struct {
	Error    bool   `json:"error"`
	Message  string `json:"message,omitempty"`
	Response any    `json:"response,omitempty"`
}

func NewApiResponse(
	isError bool,
	response any,
	message string,
) ApiResponse {
	return ApiResponse{
		Error:    isError,
		Message:  message,
		Response: response,
	}
}

func NewErrorApiResponse(message string) ApiResponse {
	return ApiResponse{
		Error:   true,
		Message: message,
	}
}

func NewSuccessApiResponse(response any) ApiResponse {
	return ApiResponse{
		Error:    false,
		Response: response,
	}
}

func ControllerWrapper(
	ctx *gin.Context,
	requestName string,
	requestType string,
	controller Controller,
) {
	db, ok := connections.GetPostgresDb()
	if !ok {
		logger.Log.WithFields(logrus.Fields{
			"controller": requestName,
		}).Error("postgres connection not found")
		ctx.JSON(http.StatusInternalServerError, NewErrorApiResponse("internal server error"))
		return
	}
	res, err := controller(ctx, db)
	if err != nil {
		var responseStatusCode int
		var responseBody ApiResponse

		if resErr, ok := err.(ExternalError); ok {
			responseStatusCode = int(resErr.Code)
			responseBody = NewErrorApiResponse(resErr.Message)
		} else {
			responseStatusCode = http.StatusInternalServerError
			responseBody = NewErrorApiResponse("internal server error")

			logger.Log.WithFields(logrus.Fields{
				"controller": requestName,
				"err":        err,
			}).Error("error in controller")
		}


		logger.Log.WithFields(logrus.Fields{
			"controller": requestName,
			"status":     responseStatusCode,
			"body":       responseBody,
		}).Info("sending response")
		ctx.JSON(responseStatusCode, responseBody)
		return
	}

	statusCode := http.StatusOK
	body := NewSuccessApiResponse(res)
	logger.Log.WithFields(logrus.Fields{
		"controller": requestName,
		"status":     statusCode,
		"body":       body,
	}).Info("sending response")
	ctx.JSON(statusCode, body)
}
