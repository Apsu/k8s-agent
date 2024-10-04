package handler

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/apsu/k8s-agent/internal/status"
	"github.com/labstack/echo/v4"
)

// streamLogs streams the logs of a systemd unit (e.g., rke2-server.service)
func Log(c echo.Context) error {
	// Use the request context to detect when the client disconnects
	ctx := c.Request().Context()

	// Run journalctl command with the request context
	cmd := exec.CommandContext(ctx, "journalctl", "-u", "rke2-"+status.GetStatus().Type, "-f")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("LOG: Failed to get stdout: %v", err))
	}

	if err := cmd.Start(); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("LOG: Failed to start journalctl: %v", err))
	}

	// Ensure the command is properly cleaned up when the client disconnects or request ends
	defer cmd.Wait()

	c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Response().WriteHeader(http.StatusOK)

	_, err = c.Response().Write([]byte("Streaming logs...\n")) // Initial message
	if err != nil {
		return err
	}

	// Stream the logs continuously
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			// Client has disconnected or request was canceled
			cmd.Process.Kill() // Ensure the command is killed
			return ctx.Err()   // Return the context error (e.g., context.Canceled)
		default:
			n, err := stdout.Read(buf)
			if err != nil {
				return err
			}
			if n > 0 {
				_, writeErr := c.Response().Write(buf[:n])
				if writeErr != nil {
					return writeErr
				}
				c.Response().Flush()
			}
		}
	}
}
