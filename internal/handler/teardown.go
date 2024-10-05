package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/apsu/k8s-agent/internal/status"
	"github.com/labstack/echo/v4"
)

func Teardown(c echo.Context) error {
	curStatus := status.GetStatus()
	if curStatus.State != status.Ready {
		return c.String(http.StatusConflict, "Agent not ready")
	}
	curStatus.State = status.Terminating
	status.UpdateStatus(curStatus)

	// Start teardown in the background
	go func() {
		fmt.Println("TEARDOWN: Starting teardown")

		cmd := exec.Command("/usr/local/bin/rke2-uninstall.sh")

		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start teardown: %v\n", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Teardown script failed: %v\n", err)
			return
		}

		curStatus.State = status.Ready
		status.UpdateStatus(curStatus)
		fmt.Println("TEARDOWN: Completed successfully")
	}()

	// Return immediately to the client
	return c.String(http.StatusAccepted, "Teardown started")
}
