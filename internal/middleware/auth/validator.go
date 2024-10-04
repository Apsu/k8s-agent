package auth

import (
	"os"

	"github.com/labstack/echo/v4"
)

var token = os.Getenv("TOKEN")

func TokenValidator(key string, c echo.Context) (bool, error) {
	return key == token, nil
}
