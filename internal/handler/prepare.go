package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/apsu/k8s-agent/internal/status"
	"github.com/labstack/echo/v4"
)

func Prepare(c echo.Context) error {
	// Start deployment in a background goroutine
	go func() {
		curStatus := status.GetStatus()
		curStatus.State = status.Preparing
		status.UpdateStatus(curStatus)

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

		fmt.Println("PREPARE: Prepare completed successfully")
	}()

	// Return immediately to the client
	return c.String(http.StatusAccepted, "Prepare started")
}
