package logger

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupLogger(e *echo.Echo) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			} else {
				fmt.Printf("REQUEST: uri: %v, status: %v, error: %v\n", v.URI, v.Status, v.Error.Error())
			}
			return nil
		},
	}))
}
