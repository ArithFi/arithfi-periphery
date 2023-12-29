package scan

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"time"
)

// FFutureTrading 扫描这个表
func FFutureTrading(c echo.Context) error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "f_future_trading_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	fmt.Println("f_future_trading last timestamp:", lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT product, positionIndex, leverage, orderType, mode, direction, margin, volume, fees, sellValue, walletAddress, kolAddress, availableBanlance, copyAccountBalance
FROM f_future_trading 
WHERE timestamp > ? 
ORDER By timestamp 
LIMIT 100
`, lastTimestamp.String())
	if err != nil {
		return err
	}
	defer query.Close()

	for query.Next() {
		var product string
		var positionIndex int64
		var timeStamp int
		var leverage int64
		var orderType string
		var mode string
		var direction string
		var margin float64
		var volume float64
		var fees float64
		var sellValue float64
		var walletAddress string
		var kolAddress string
		var availableBanlance float64
		var copyAccountBalance float64

		err := query.Scan(&product, &positionIndex, &leverage, &orderType, &mode, &direction, &margin, &volume, &fees, &sellValue, &walletAddress, &kolAddress, &availableBanlance, &copyAccountBalance)
		if err != nil {
			return err
		}

		fmt.Println("f_future_trading:", product, positionIndex, leverage, orderType, mode, direction, margin, volume, fees, sellValue, walletAddress, kolAddress, availableBanlance, copyAccountBalance)

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timeStamp)+8*60*60, 0).Format("2006-01-02")
		fmt.Println("date:", date)
		if orderType == "MARKET_ORDER_FEE" || orderType == "LIMIT_ORDER_FEE" {
			// 处理每个用户每天开单的数据汇总，新增仓位，开单数量
			handleNewOrder(mode)
		} else if orderType == "MARKET_CLOSE_FEE" || orderType == "TP_ORDER_FEE" || orderType == "SL_ORDER_FEE" || orderType == "MARKET_LIQUIDATION" {
			// 处理每个每天的净销毁
			handleBurn(mode)
		}
	}

	return nil
}

func handleNewOrder(mode string) {

}

func handleBurn(mode string) {

}
