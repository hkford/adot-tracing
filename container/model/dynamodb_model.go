package model

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBTableAPI interface {
	DescribeTable(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)

	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type TableInfo struct {
	Status string `json:"status"`
	Items  int64  `json:"items"`
}

func GetTableInfo(c context.Context, api DynamoDBTableAPI, input *dynamodb.DescribeTableInput) (*TableInfo, error) {
	resp, err := api.DescribeTable(c, input)
	if err != nil {
		log.Error().Err(err).Msg("Error response calling dynamodb.DescribeTable")
		return nil, err
	}
	info := &TableInfo{
		Status: string(resp.Table.TableStatus),
		Items:  resp.Table.ItemCount,
	}
	return info, nil
}

type Movie struct {
	ID    int    `dynamodbav:"id" validate:"required"`
	Title string `dynamodbav:"title" validate:"required"`
}

func CreateMovie(c context.Context, api DynamoDBTableAPI, input *dynamodb.PutItemInput) error {
	_, err := api.PutItem(c, input)
	if err != nil {
		log.Error().Err(err).Msg("Error response calling dynamodb.PutItem")
		return err
	}
	return nil
}

type MovieWithID struct {
	ID int `dynamodbav:"id" validate:"required"`
}

func ReadMovie(c context.Context, api DynamoDBTableAPI, input *dynamodb.GetItemInput) (*Movie, error) {
	output, err := api.GetItem(c, input)
	if err != nil {
		log.Error().Err(err).Msg("Error response calling dynamodb.GetItem")
		return nil, err
	}
	var movie Movie
	err = attributevalue.UnmarshalMap(output.Item, &movie)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal new movie item")
		return nil, err
	}
	if len(output.Item) == 0 {
		return nil, nil
	}
	return &movie, nil
}

func DeleteMovie(c context.Context, api DynamoDBTableAPI, input *dynamodb.DeleteItemInput) error {
	_, err := api.DeleteItem(c, input)
	if err != nil {
		log.Error().Err(err).Msg("Error response calling dynamodb.DeleteItem")
		return err
	}
	return nil
}
