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

	fmt.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT product, positionIndex, leverage, orderType, mode, direction, margin, volume, sellValue, walletAddress, kolAddress, availableBanlance, copyAccountBalance
FROM f_future_trading 
WHERE timestamp > ? 
ORDER By timestamp 
LIMIT 100
`, lastTimestamp.Val())
	if err != nil {
		return err
	}
	defer query.Close()
	var newLastTimestamp int
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
		var sellValue float64
		var walletAddress string
		var kolAddress string
		var availableBanlance float64
		var copyAccountBalance float64

		err := query.Scan(&product, &positionIndex, &leverage, &orderType, &mode, &direction, &margin, &volume, &sellValue, &walletAddress, &kolAddress, &availableBanlance, &copyAccountBalance)
		if err != nil {
			return err
		}
		newLastTimestamp = timeStamp

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timeStamp)+8*60*60, 0).Format("2006-01-02")
		fmt.Println("date:", date)
		if orderType == "MARKET_ORDER_FEE" || orderType == "LIMIT_ORDER_FEE" {
			// 处理每个用户每天开单的数据汇总，新增仓位，开单数量
			handleNewOrder(mode, date, walletAddress, kolAddress, volume)
		} else if orderType == "MARKET_CLOSE_FEE" || orderType == "TP_ORDER_FEE" || orderType == "SL_ORDER_FEE" || orderType == "MARKET_LIQUIDATION" {
			// 处理每个每天的净销毁
			handleBurn(mode, sellValue, margin, walletAddress, kolAddress, date)
		}
	}
	cache.CACHE.Set(ctx, "f_future_trading_last_timestamp", newLastTimestamp, 0)
	return nil
}

func handleNewOrder(mode string, date string, walletAddress string, kolAddress string, volume float64) {
	_, err := mysql.MYSQL.Exec(`
INSERT INTO b_daily_offchain_futures_metrics (date, walletAddress, mode, kolAddress, new_position_counts, new_position_size)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE new_position_counts = new_position_counts + 1, new_position_size = new_position_size + ?
`, date, walletAddress, mode, kolAddress, 1, volume, volume)
	if err != nil {
		fmt.Println("handleNewOrder err:", err)
		return
	}
	fmt.Println("handleNewOrder ok")
}

func handleBurn(mode string, sellValue float64, margin float64, walletAddress string, kolAddress string, date string) {
	var netBurnAmount = sellValue - margin
	_, err := mysql.MYSQL.Exec(`
INSERT INTO b_daily_offchain_futures_metrics (date, walletAddress, mode, kolAddress, net_burn_amount)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE net_burn_amount = VALUES(net_burn_amount)
`, date, walletAddress, mode, kolAddress, netBurnAmount)
	if err != nil {
		fmt.Println("handleBurn err:", err)
		return
	}
	fmt.Println("handleBurn ok")
}
