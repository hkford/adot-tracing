package server

import (
	"fmt"
	"main/conf"
	"main/controller"
	"main/util"
	"os"
)

func Run() {

	conf.RegisterTracerProvider()

	dynamodbClient := conf.NewDynamoDBClient()
	dynamodbClientConfig := conf.NewDynamodbClientConfig()

	tableController := controller.NewTableController(dynamodbClient, *dynamodbClientConfig)
	movieController := controller.NewMovieController(dynamodbClient, *dynamodbClientConfig)
	baseController := controller.NewBaseController()

	c := Controllers{
		baseController,
		tableController,
		movieController,
	}

	e := NewRouter(c)

	e = RegisterMiddlewares(e)

	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		util.FatalLog(nil, "APP_PORT not set")
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
