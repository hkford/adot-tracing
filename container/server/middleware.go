package server

import (
	"main/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func RegisterMiddlewares(e *echo.Echo) *echo.Echo {
	e.Use(otelecho.Middleware("my-server"))
	e.Use(middleware.RequestID())
	e.Use(util.WrapEchoHandlerWithLogging)
	return e
}
