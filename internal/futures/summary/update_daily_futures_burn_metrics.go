package summary

import (
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	UpdateDailyFuturesDestroyMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date"` // Symbol
	}
)

/*
UpdateDailyFuturesBurnMetrics UpdateDailyFuturesBurnMetrics
Weight: 1
Parameters: NONE
*/
func UpdateDailyFuturesBurnMetrics(c echo.Context) error {
	// get date from query
	var req UpdateDailyDepositWithdrawMetricsReqType
	if err := c.Bind(&req); err != nil {
		return err
	}
	date := req.Date

	from := date + " 00:00:00"
	to := date + " 23:59:59"

	handleFuturesBurn(from, to, date)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}

func handleFuturesBurn(from string, to string, date string) {
	rows, err := mysql.MYSQL.Query(`SELECT walletAddress, kolAddress, mode, SUM(SellValue - Margin) FROM f_future_trading
WHERE CONVERT_TZ(timeStamp, '+00:00', '+08:00') >= ?
  AND CONVERT_TZ(timeStamp, '+00:00', '+08:00') <= ?
  AND orderType in ('MARKET_CLOSE_FEE', 'SL_ORDER_FEE', 'TP_ORDER_FEE', 'MARKET_LIQUIDATION')
GROUP BY walletAddress, kolAddress, mode;`, from, to)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var walletAddress, kolAddress, mode string
		var net_burn_amount float64
		err = rows.Scan(&walletAddress, &kolAddress, &mode, &net_burn_amount)
		if err != nil {
			return
		}
		fmt.Println(walletAddress, kolAddress, mode, net_burn_amount)

		_, err2 := mysql.MYSQL.Exec(`INSERT INTO b_daily_offchain_futures_metrics (walletAddress, kolAddress, mode, date, net_burn_amount)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE net_burn_amount = VALUES(net_burn_amount);`, walletAddress, kolAddress, mode, date, net_burn_amount)
		if err2 != nil {
			return
		}
	}
}
