package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Hello function to return Hello, World!
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
