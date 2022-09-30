package controller

import (
	"fmt"
	"main/conf"
	"main/model"
	"main/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MovieController interface {
	Create(c echo.Context) error
	Read(c echo.Context) error
	Delete(c echo.Context) error
}

type movieController struct {
	client model.DynamoDBTableAPI
	config conf.DynamodbClientConfig
}

func NewMovieController(dynamodbClient model.DynamoDBTableAPI, dynamodbClientConfig conf.DynamodbClientConfig) MovieController {
	return &movieController{dynamodbClient, dynamodbClientConfig}
}

func (mc *movieController) Create(c echo.Context) error {
	tracer := otel.Tracer("otel-tracer")
	ctx, span := tracer.Start(c.Request().Context(), "POST /movie")
	defer span.End()
	m := new(model.Movie)
	if err := c.Bind(m); err != nil {
		msg := fmt.Sprintf("Error binding request body to movie %#v", m)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		msg := fmt.Sprintf("Error validating request body to movie %#v", m)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	movie, err := attributevalue.MarshalMap(m)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling new movie item %#v", m)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      movie,
		TableName: &mc.config.TableName,
	}

	err = model.CreateMovie(ctx, mc.client, input)
	if err == nil {
		msg := "Created new movie"
		util.ResponseLog(msg)
		response := &jsonResponse{
			Message: msg,
		}
		return c.JSON(http.StatusOK, response)
	} else {
		util.ErrorLog(err, "Failed to create new movie")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
}

func (mc *movieController) Read(c echo.Context) error {
	tracer := otel.Tracer("otel-tracer")
	ctx, span := tracer.Start(c.Request().Context(), "GET /movie/:id")
	defer span.End()
	idFromRequest := c.Param("id")
	id, err := strconv.Atoi(idFromRequest)
	if err != nil {
		msg := fmt.Sprintf("Error converting id to integer: %s", idFromRequest)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	m := model.MovieWithID{
		ID: id,
	}
	movieInput, err := attributevalue.MarshalMap(m)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling new movie item %#v", m)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	input := &dynamodb.GetItemInput{
		Key:       movieInput,
		TableName: &mc.config.TableName,
	}

	movie, err := model.ReadMovie(ctx, mc.client, input)
	if err != nil {
		util.ErrorLog(err, fmt.Sprintf("Failed to get movie for id: %s", idFromRequest))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if movie != nil {
		return c.JSON(http.StatusOK, movie)
	} else {
		msg := fmt.Sprintf("No movie found for id: %s", idFromRequest)
		util.ResponseLog(msg)
		response := &jsonResponse{
			Message: msg,
		}
		return c.JSON(http.StatusOK, response)
	}
}

func (mc *movieController) Delete(c echo.Context) error {
	tracer := otel.Tracer("otel-tracer")
	ctx, span := tracer.Start(c.Request().Context(), "DELETE /movie/:id")
	defer span.End()

	idFromRequest := c.Param("id")
	id, err := strconv.Atoi(idFromRequest)
	if err != nil {
		msg := fmt.Sprintf("Error converting id to integer: %s", idFromRequest)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	m := model.MovieWithID{
		ID: id,
	}
	movieInput, err := attributevalue.MarshalMap(m)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling new movie item %#v", m)
		util.ErrorLog(err, msg)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	input := &dynamodb.DeleteItemInput{
		Key:       movieInput,
		TableName: &mc.config.TableName,
	}

	err = model.DeleteMovie(ctx, mc.client, input)
	if err != nil {
		util.ErrorLog(err, fmt.Sprintf("Failed to delete movie for id: %s", idFromRequest))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else {
		msg := fmt.Sprintf("Deleted movie for id: %s", idFromRequest)
		util.ResponseLog(msg)
		response := &jsonResponse{
			Message: msg,
		}
		return c.JSON(http.StatusOK, response)
	}
}
