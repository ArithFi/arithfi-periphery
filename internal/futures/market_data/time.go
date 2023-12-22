package market_data

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

/*
Time Test connectivity to the Rest API and get the current server time.
Weight: 1
Parameters: NONE
*/
func Time(c echo.Context) error {
	now := time.Now().Second()

	return c.JSON(http.StatusOK, map[string]int{
		"serverTime": now,
	})

}
