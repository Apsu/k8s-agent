package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/apsu/k8s-agent/internal/status"
	"github.com/labstack/echo/v4"
)

func Prepare(c echo.Context) error {
	curStatus := status.GetStatus()
	if curStatus.State != status.Running {
		return c.String(http.StatusConflict, "Agent not ready")
	}
	curStatus.State = status.Preparing
	status.UpdateStatus(curStatus)

	// Start prepare in the background
	go func() {

		fmt.Println("PREPARE: Starting prepare node")

		cmd := exec.Command("/opt/k8s-agent/scripts/prepare.sh")

		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start prepare: %v\n", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Prepare script failed: %v\n", err)
			return
		}

		curStatus.State = status.Running
		status.UpdateStatus(curStatus)
		fmt.Println("PREPARE: Prepare completed successfully")
	}()

	// Return immediately to the client
	return c.String(http.StatusAccepted, "Prepare started")
}
