package util

import (
	"fmt"

	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
)

func FatalLog(err error, message string) {
	const fatalLevel slog.Level = 10
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Log(context.Background(), fatalLevel, fmt.Sprintf("%s: %v", message, err))
	os.Exit(1)
}

func ErrorLog(err error, message string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Error(fmt.Sprintf("%s: %v", message, err))
}

type requestLogInput struct {
	method     string
	path       string
	clientAddr string
	requestID  string
}

func requestLog(i requestLogInput) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("access log", slog.String("method", i.method), slog.String("path", i.path), slog.String("client", i.clientAddr), slog.String("request id", i.requestID))
}

func ResponseLog(msg string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(msg)
}

func WrapEchoHandlerWithLogging(originalHandler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		i := requestLogInput{
			method:     c.Request().Method,
			path:       c.Request().URL.Path,
			clientAddr: c.Request().RemoteAddr,
			requestID:  c.Response().Header().Get(echo.HeaderXRequestID),
		}
		requestLog(i)
		return originalHandler(c)
	}
}
