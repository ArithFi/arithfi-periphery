package summary

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// 每日的期货数据，用户、用户类型、交易类型、KOL、新增交易笔数、新增交易规模、净销毁数量

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
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// 需要连接数据库，获取 f_future_trading 当天的信息

	db := mysql.MYSQL
	// 获取 walletAddress，kolAddress，new_position_counts, mode, new_position_size
	rows, err := db.Query(`SELECT walletAddress, kolAddress, COUNT(positionIndex), mode, SUM( margin * leverage) positionSize FROM f_future_trading
WHERE CONVERT_TZ(timeStamp, '+00:00', '+08:00') >= '2023-12-01 00:00:00'
  AND CONVERT_TZ(timeStamp, '+00:00', '+08:00') < '2023-12-02 00:00:00'
  AND orderType in ("MARKET_ORDER_FEE", "LIMIT_ORDER_FEE")
GROUP BY walletAddress, kolAddress, mode`, date, date)
	if err != nil {
		return err
	}
	defer rows.Close()
	// 获取到信息后，遍历并更新数据

	return c.JSON(http.StatusOK, map[string]string{})
}
