package util

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func PanicLog(err error, message string) {
	if err != nil {
		log.Panic().Msg(fmt.Sprintf("%s: %v", message, err))
	}
}

func ErrorLog(err error, message string) {
	log.Error().Err(err).Msg(message)
}

type requestLogInput struct {
	method     string
	path       string
	clientAddr string
	requestID  string
}

func requestLog(i requestLogInput) {
	log.Info().Str("method", i.method).Str("path", i.path).Str("client", i.clientAddr).Str("request id", i.requestID).Send()
}

func ResponseLog(msg string) {
	log.Info().Str("message", msg).Send()
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
