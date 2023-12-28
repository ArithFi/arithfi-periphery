package summary

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
)

// 每日的充值和提现汇总数据
// 用户、总充值、总提现

type (
	UpdateDailyDepositWithdrawMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date" validate:"required"` // Symbol
	}
)

/*
UpdateDailyDepositWithdrawMetrics UpdateDailyDepositWithdrawMetrics
Weight: 1
Parameters: NONE
*/
func UpdateDailyDepositWithdrawMetrics(c echo.Context) error {
	// get date from query
	var req UpdateDailyDepositWithdrawMetricsReqType
	if err := c.Bind(&req); err != nil {
		return err
	}
	date := req.Date

	from := date + " 00:00:00"
	to := date + " 23:59:59"

	handleWithdraw(from, to, date)
	handleDeposit(from, to, date)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "success",
		"date":    date,
	})
}

func handleWithdraw(from string, to string, date string) {
	withdrawRows, err := mysql.MYSQL.Query(`SELECT walletAddress, SUM(amount), COUNT(ordertype)
  FROM deposit_withdrawal
  WHERE CONVERT_TZ(_createTime, '+00:00', '+08:00') >= ?
    AND CONVERT_TZ(_createTime, '+00:00', '+08:00') <= ?
    AND orderType IN ('WITHDRAW')
  GROUP BY walletAddress`, from, to)
	if err != nil {
		return
	}
	defer withdrawRows.Close()
	for withdrawRows.Next() {
		var walletAddress string
		var withdrawAmount float64
		var withdrawCounts int64
		if err := withdrawRows.Scan(&walletAddress, &withdrawAmount, &withdrawCounts); err != nil {
			return
		}

		_, err := mysql.MYSQL.Exec(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (wallet, date, withdraw_amount, withdraw_counts)
				VALUES (?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE
				withdraw_amount = withdraw_amount,
				withdraw_counts = withdraw_counts`, walletAddress, date, withdrawAmount, withdrawCounts)
		if err != nil {
			return
		}
	}
}

func handleDeposit(from string, to string, date string) {
	depositRows, err := mysql.MYSQL.Query(`SELECT walletAddress, SUM(amount), COUNT(ordertype)
  FROM deposit_withdrawal
  WHERE CONVERT_TZ(_createTime, '+00:00', '+08:00') >= ?
    AND CONVERT_TZ(_createTime, '+00:00', '+08:00') <= ?
    AND orderType IN ('DEPOSIT', 'WALLET_DEPOSIT')
  GROUP BY walletAddress`, from, to)
	if err != nil {
		return
	}
	defer depositRows.Close()
	for depositRows.Next() {
		var walletAddress string
		var depositAmount float64
		var depositCounts int64
		if err := depositRows.Scan(&walletAddress, &depositAmount, &depositCounts); err != nil {
			return
		}

		_, err := mysql.MYSQL.Exec(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (wallet, date, deposit_amount, deposit_counts)
				VALUES (?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE
				deposit_amount = deposit_amount,
				deposit_counts = deposit_counts`, walletAddress, date, depositAmount, depositCounts)
		if err != nil {
			return
		}
	}
}
