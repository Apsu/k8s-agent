package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/apsu/k8s-agent/internal/status"
)

func Status(c echo.Context) error {
	// Placeholder for status
	return c.JSON(http.StatusOK, status.GetStatus())
}
