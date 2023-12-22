package market_data

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
Ping Test connectivity to the Rest API.
Weight: 1
Parameters: NONE
*/
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}
