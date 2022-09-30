package server

import (
	controller "main/controller"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	bc controller.BaseController
	tc controller.TableController
	mc controller.MovieController
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewRouter(c Controllers) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.GET("/", c.bc.HealthCheck)
	e.GET("/traceid", c.bc.GetTraceID)
	e.GET("/table", c.tc.Describe)

	e.POST("/movie", c.mc.Create)
	e.GET("/movie/:id", c.mc.Read)
	e.DELETE("/movie/:id", c.mc.Delete)

	return e
}
