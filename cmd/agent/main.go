package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/apsu/k8s-agent/internal/handler"
	"github.com/apsu/k8s-agent/internal/middleware/auth"
	"github.com/apsu/k8s-agent/internal/middleware/logger"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	logger.SetupLogger(e)
	e.Use(middleware.Recover())

	// Public routes
	e.GET("/status", handler.Status)
	e.GET("/log", handler.Log)

	// Setup authenticated route group
	g := e.Group("/api")

	// Token-based auth
	g.Use(middleware.KeyAuth(auth.TokenValidator))

	// Admin routes
	g.POST("/api/deploy", handler.Deploy)
	g.POST("/teardown", handler.Teardown)
	g.POST("/restart", handler.Restart)
	g.POST("/upgrade", handler.Upgrade)

	// Start server
	e.Logger.Fatal(e.Start(":8765"))
}
