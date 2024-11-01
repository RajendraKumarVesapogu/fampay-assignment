package routes

import (
	"fampay-assignment/middleware"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	e := gin.New()

	e.Use(middleware.Timeout())
	e.Use(middleware.Cors())
	e.Use(
		gin.LoggerWithWriter(gin.DefaultWriter),
		gin.Recovery(),
	)
	e.Use(middleware.RequestLogger())

	return e
}
