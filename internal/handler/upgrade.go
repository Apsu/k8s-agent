package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	// "github.com/apsu/k8s-agent/internal/status"
)

func Upgrade(c echo.Context) error {
	// Placeholder for restart
	return echo.NewHTTPError(http.StatusNotImplemented)
}
