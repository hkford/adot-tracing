package conf

import (
	"context"
	"log"
	"main/model"
	"main/util"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type DynamodbClientConfig struct {
	TableName string
}

func NewDynamodbClientConfig() *DynamodbClientConfig {
	table := os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("No TABLE_NAME set")
	}
	config := &DynamodbClientConfig{table}
	return config
}

func NewDynamoDBClient() model.DynamoDBTableAPI {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		util.PanicLog(err, "unable to load SDK config")
	}
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	client := dynamodb.NewFromConfig(cfg)
	return client
}
