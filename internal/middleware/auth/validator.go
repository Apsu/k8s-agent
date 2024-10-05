package auth

import (
	"os"

	"github.com/labstack/echo/v4"
)

func TokenValidator(key string, c echo.Context) (bool, error) {
	token := os.Getenv("TOKEN")
	return key == token, nil
}
