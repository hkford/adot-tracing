package controller

import (
	"errors"
	"fmt"
	"main/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/trace"
)

type BaseController interface {
	HealthCheck(c echo.Context) error
	GetTraceID(c echo.Context) error
}

type baseController struct {
	idGenerator *xray.IDGenerator
}

type traceIDsResponse struct {
	OtelTraceIDFromSpanContext string
	XRayTraceIDFromSpanContext string
	OtelTraceIDFromGenerator   string
	XRayTraceIDFromGenerator   string
}

func NewBaseController() BaseController {
	idGenerator := xray.NewIDGenerator()
	return &baseController{idGenerator}
}

type jsonResponse struct {
	Message string
}

func (bc *baseController) HealthCheck(c echo.Context) error {
	response := &jsonResponse{
		Message: "Healthy",
	}
	return c.JSON(http.StatusOK, response)
}

func convertTraceFromOtelToXRay(otelTraceID string) (string, error) {
	if len(otelTraceID) == 32 {
		xrayTraceID := fmt.Sprintf("1-%s-%s", otelTraceID[0:8], otelTraceID[8:])
		return xrayTraceID, nil
	} else {
		message := fmt.Sprintf("Cannot convert trace id to X-Ray format from %s", otelTraceID)
		return "", errors.New(message)
	}
}

func (bc *baseController) GetTraceID(c echo.Context) error {
	otelTraceIDFromSpanContext := trace.SpanFromContext(c.Request().Context()).SpanContext().TraceID().String()
	xrayTraceIDFromSpanContext, err := convertTraceFromOtelToXRay(otelTraceIDFromSpanContext)
	if err != nil {
		util.ErrorLog(err, "Error converting traceid from span context")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	otelTraceIDFromGenerator, _ := bc.idGenerator.NewIDs(c.Request().Context())
	xrayTraceIDFromGenerator, err := convertTraceFromOtelToXRay(otelTraceIDFromGenerator.String())
	if err != nil {
		util.ErrorLog(err, "Error converting traceid from generator")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := &traceIDsResponse{
		OtelTraceIDFromSpanContext: otelTraceIDFromSpanContext,
		XRayTraceIDFromSpanContext: xrayTraceIDFromSpanContext,
		OtelTraceIDFromGenerator:   otelTraceIDFromGenerator.String(),
		XRayTraceIDFromGenerator:   xrayTraceIDFromGenerator,
	}
	return c.JSON(http.StatusOK, response)
}
