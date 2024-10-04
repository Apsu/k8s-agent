package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	// "github.com/apsu/k8s-agent/internal/status"
)

func Teardown(c echo.Context) error {
	// Placeholder for teardown
	return echo.NewHTTPError(http.StatusNotImplemented)
}
