package summary

import (
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type (
	UpdateDailyTradeMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date" validate:"required"` // Symbol
	}
)

func UpdateDailyTradeMetrics(c echo.Context) error {
	// get date from query
	var req UpdateDailyTradeMetricsReqType
	if err := c.Bind(&req); err != nil {
		return err
	}
	date := req.Date

	from := date + " 00:00:00"
	to := date + " 23:59:59"

	handleBuy(from, to, date)
	handleSell(from, to, date)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}

func handleBuy(from string, to string, date string) {
	query := `
	SELECT to_address, SUM(value), COUNT(to_address)
	FROM erc20_transfer_bsc
	WHERE from_address = ?
	AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') >= ?
	AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') <= ?
	GROUP BY to_address;
	`
	rows, err := mysql.MYSQL.Query(query, from, to, date)
	if err != nil {
		return
	}
	defer rows.Close()

	var values []string
	var args []interface{}
	insertQuery := `INSERT INTO b_daily_onchain_trade_metrics (date, walletAddress, buy_amount, buy_counts)
	VALUES %s ON DUPLICATE KEY UPDATE buy_amount = VALUES(buy_amount), buy_counts = VALUES(buy_counts);`

	for rows.Next() {
		var walletAddress string
		var buyAmount, buyCounts float64
		if err := rows.Scan(&walletAddress, &buyAmount, &buyCounts); err != nil {
			return
		}

		values = append(values, "(?, ?, ?, ?)")
		args = append(args, date, walletAddress, buyAmount, buyCounts)
	}

	if len(values) == 0 {
		return
	}

	insertQuery = fmt.Sprintf(insertQuery, strings.Join(values, ","))
	_, err = mysql.MYSQL.Exec(insertQuery, args...)
	if err != nil {
		return
	}

	return
}

func handleSell(from string, to string, date string) {
	query := `SELECT from_address, SUM(value), COUNT(from_address)
  FROM erc20_transfer_bsc
  WHERE to_address = '0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38'
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') >= ?
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') <= ?
  GROUP BY from_address;`

	rows, err := mysql.MYSQL.Query(query, from, to)
	if err != nil {
		return
	}

	var values []string
	var args []interface{}
	insertQuery := `INSERT INTO b_daily_onchain_trade_metrics (date, walletAddress, sell_amount, sell_counts)
	VALUES %s ON DUPLICATE KEY UPDATE sell_amount = VALUES(sell_amount), sell_counts = VALUES(sell_counts);`

	for rows.Next() {
		var walletAddress string
		var sellAmount, sellCounts float64
		if err := rows.Scan(&walletAddress, &sellAmount, &sellCounts); err != nil {
			return
		}

		values = append(values, "(?, ?, ?, ?)")
		args = append(args, date, walletAddress, sellAmount, sellCounts)
	}

	if len(values) == 0 {
		return
	}

	insertQuery = fmt.Sprintf(insertQuery, strings.Join(values, ","))
	_, err = mysql.MYSQL.Exec(insertQuery, args...)
	if err != nil {
		return
	}
	return
}
