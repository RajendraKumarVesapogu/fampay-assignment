package routes

import (
	"fampay-assignment/middleware"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	e := gin.New()

	e.Use(middleware.Cors())
	e.Use(middleware.Timeout())
	e.Use(
		gin.LoggerWithWriter(gin.DefaultWriter),
		gin.Recovery(),
	)
	e.Use(middleware.RequestLogger())

	Videos(e)
	
	return e
}
