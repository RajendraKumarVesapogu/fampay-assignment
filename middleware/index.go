package middleware

import (
	"bytes"
	"io"
	"net/http"

	"fampay-assignment/config"
	"fampay-assignment/lib"
	"fampay-assignment/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		
		query := ctx.Request.URL.Query()
		body, _ := io.ReadAll(ctx.Request.Body)

		logger.Log.WithFields(logger.Fields{
			"path":    path,
			"query":   query,
			"body":    string(body),
		}).Info("request received")

		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		ctx.Next()
	}
}


func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowCredentials: true,
		AllowWildcard:    true,
	})
}

func Timeout() gin.HandlerFunc {
	return timeout.Timeout(
		timeout.WithErrorHttpCode(http.StatusRequestTimeout),
		timeout.WithDefaultMsg(lib.NewErrorApiResponse("request timeout")),
		timeout.WithGinCtxCallBack(func(ctx *gin.Context) {
			logger.Log.WithFields(logger.Fields{
				"url": ctx.Request.URL.String(),
			}).Error("request timeout")
			ctx.Abort()
		}))
}