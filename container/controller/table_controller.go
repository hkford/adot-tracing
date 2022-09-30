package controller

import (
	"main/conf"
	"main/model"
	"main/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TableController interface {
	Describe(c echo.Context) error
}

type tableController struct {
	client model.DynamoDBTableAPI
	config conf.DynamodbClientConfig
}

func NewTableController(dynamodbClient model.DynamoDBTableAPI, dynamodbClientConfig conf.DynamodbClientConfig) TableController {
	return &tableController{dynamodbClient, dynamodbClientConfig}
}

func (tc *tableController) Describe(c echo.Context) error {
	tracer := otel.Tracer("otel-tracer")
	ctx, span := tracer.Start(c.Request().Context(), "GET /table")

	span.SetAttributes(attribute.String("AWSService", "DynamoDB"))
	defer span.End()

	input := &dynamodb.DescribeTableInput{
		TableName: &tc.config.TableName,
	}

	info, err := model.GetTableInfo(ctx, tc.client, input)
	if err != nil {
		msg := "Failed to describe table"
		util.ErrorLog(err, msg)
		response := &jsonResponse{
			Message: msg,
		}
		return c.JSON(http.StatusInternalServerError, response)
	} else {
		util.ResponseLog("Described table")
		return c.JSON(http.StatusOK, info)
	}
}
