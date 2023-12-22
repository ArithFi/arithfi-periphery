package market_data

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

/*
Kline /candlestick bars for a symbol.
Kline are uniquely identified by their open time.

Weight: based on parameter LIMIT

	[1,100)	1
	[100, 500)	2
	[500, 1000]	5
	> 1000	10

Parameters:

	symbol	STRING	YES
	interval	ENUM	YES
	startTime	LONG	NO
	endTime	LONG	NO
	limit	INT	NO
*/
func Kline(c echo.Context) error {
	now := time.Now().Second()

	return c.JSON(http.StatusOK, map[string]int{
		"serverTime": now,
	})
}
