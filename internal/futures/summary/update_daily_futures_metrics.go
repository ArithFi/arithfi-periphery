package summary

import (
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	UpdateDailyMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date"` // Symbol
	}
)

func UpdateDailyFuturesMetrics(c echo.Context) error {
	// get date from query
	var req UpdateDailyMetricsReqType
	if err := c.Bind(&req); err != nil {
		return err
	}
	date := req.Date

	from := date + " 00:00:00"
	to := date + " 23:59:59"

	handleMetrics(from, to, date)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}

func handleMetrics(from string, to string, date string) {
	rows, err := mysql.MYSQL.Query(`SELECT walletAddress, kolAddress, COUNT(positionIndex), mode, SUM( margin * leverage) FROM f_future_trading
WHERE CONVERT_TZ(timeStamp, '+00:00', '+08:00') >= ?
  AND CONVERT_TZ(timeStamp, '+00:00', '+08:00') <= ?
  AND orderType in ('MARKET_ORDER_FEE', 'LIMIT_ORDER_FEE')
GROUP BY walletAddress, kolAddress, mode;`, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var walletAddress, kolAddress, mode string
		var newPositionSize float64
		var newPositionCounts int
		if err := rows.Scan(&walletAddress, &kolAddress, &newPositionCounts, &mode, &newPositionSize); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(walletAddress, kolAddress, mode, newPositionCounts, newPositionSize)

		_, err := mysql.MYSQL.Exec(`INSERT INTO b_daily_offchain_futures_metrics (walletAddress, kolAddress, mode, date, new_position_counts, new_position_size)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE new_position_counts = VALUES(new_position_counts), new_position_size = VALUES(new_position_size);`, walletAddress, kolAddress, mode, date, newPositionCounts, newPositionSize)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
