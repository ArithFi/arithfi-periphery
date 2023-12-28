package summary

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
SELECT walletAddress, kolAddress, mode, SUM(SellValue - Margin) FROM f_future_trading
WHERE CONVERT_TZ(timeStamp, '+00:00', '+08:00') >= '2023-12-01 00:00:00'
  AND CONVERT_TZ(timeStamp, '+00:00', '+08:00') < '2023-12-02 00:00:00'
  AND orderType in ("MARKET_CLOSE_FEE", "SL_ORDER_FEE", "TP_ORDER_FEE", "MARKET_LIQUIDATION")
GROUP BY walletAddress, kolAddress, mode;
*/

type (
	UpdateDailyFuturesDestroyMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date"` // Symbol
	}
)

/*
UpdateDailyFuturesDestroyMetrics UpdateDailyFuturesDestroyMetrics
Weight: 1
Parameters: NONE
*/
func UpdateDailyFuturesDestroyMetrics(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}
