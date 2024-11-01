package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"fampay-assignment/config"
	"fampay-assignment/logger"
	"fampay-assignment/routes"
)

var router *gin.Engine
var ginLambda *ginadapter.GinLambda

func init() {
	logger.Log.Info("starting the server")
	router = routes.Router()
	ginLambda = ginadapter.New(router)
}

func GinRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	port := fmt.Sprintf(":%s", config.Port)
	logger.Log.WithFields(logger.Fields{
		"port":       port,
		"GOMAXPROCS": runtime.GOMAXPROCS(0),
	}).Info("server started listening..")

	router.Run(port)

}
