package market_data

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
ExchangeInfo Current exchange trading rules and symbol information
Weight: 1
Parameters: NONE
*/
func ExchangeInfo(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
