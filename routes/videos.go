package routes

import (
	"fampay-assignment/controllers"
	"fampay-assignment/lib"

	"github.com/gin-gonic/gin"
)

func Videos(engine *gin.Engine) *gin.Engine {
	videos := engine.Group("/videos")

	videos.GET("/", func(ctx *gin.Context) {
		lib.ControllerWrapper(ctx, "GetLatestVideos", controllers.GetLatestVideos)
	})

	videos.POST("/key", func(ctx *gin.Context) {
		lib.ControllerWrapper(ctx, "AddYoutubeAPIKey", controllers.AddYoutubeAPIKey)
	})

	return engine
}
