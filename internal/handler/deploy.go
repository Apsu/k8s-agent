package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/apsu/k8s-agent/internal/status"
	"github.com/labstack/echo/v4"
)

func Deploy(c echo.Context) error {
	curStatus := status.GetStatus()
	if curStatus.State != status.Ready {
		return c.String(http.StatusConflict, "Agent not ready")
	}
	curStatus.State = status.Deploying
	status.UpdateStatus(curStatus)

	// Start deployment in a background goroutine
	go func() {
		fmt.Println("DEPLOY: Starting deployment")

		cmd := exec.Command("/opt/k8s-agent/scripts/deploy.sh")

		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start deployment: %v\n", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Deployment script failed: %v\n", err)
			return
		}

		curStatus.State = status.Ready
		status.UpdateStatus(curStatus)
		fmt.Println("DEPLOY: Deployment completed successfully")
	}()

	// Return immediately to the client
	return c.String(http.StatusAccepted, "Deployment started")
}
