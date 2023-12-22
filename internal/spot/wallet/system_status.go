package wallet

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
GetSystemStatus Get system status
*/
func GetSystemStatus(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
	})
}
