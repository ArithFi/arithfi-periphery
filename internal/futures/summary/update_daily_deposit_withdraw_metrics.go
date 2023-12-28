package summary

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// 每日的充值和提现汇总数据
// 用户、总充值、总提现
/*
  SELECT walletAddress, SUM(amount) AS deposit
  FROM deposit_withdrawal
  WHERE CONVERT_TZ(_createTime, '+00:00', '+08:00') >= '2023-12-01 00:00:00'
    AND CONVERT_TZ(_createTime, '+00:00', '+08:00') < '2023-12-02 00:00:00'
    AND orderType IN ('DEPOSIT', 'WALLET_DEPOSIT')
  GROUP BY walletAddress
*/

type (
	UpdateDailyDepositWithdrawMetricsReqType struct {
		Date string `query:"date" form:"date" json:"date"` // Symbol
	}
)

/*
UpdateDailyDepositWithdrawMetrics UpdateDailyDepositWithdrawMetrics
Weight: 1
Parameters: NONE
*/
func UpdateDailyDepositWithdrawMetrics(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}
