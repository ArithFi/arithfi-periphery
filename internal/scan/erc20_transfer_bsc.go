package scan

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

// ERC20TransferBSC 扫描这个表的事件
func ERC20TransferBSC(c echo.Context) error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "erc20_transfer_bsc_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	fmt.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT from_address, to_address, timestamp, value
FROM erc20_transfer_bsc 
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
		var fromAddress, toAddress string
		var timestamp int
		var value float64
		err := query.Scan(&fromAddress, &toAddress, &timestamp, &value)
		if err != nil {
			return err
		}
		fmt.Println("fromAddress:", fromAddress)
		fmt.Println("toAddress:", toAddress)
		fmt.Println("timestamp:", timestamp)
		newLastTimestamp = timestamp
		fmt.Println("value:", value)

		newFromBalance := cache.CACHE.IncrByFloat(ctx, "BALANCE#"+fromAddress, -value)
		newToBalance := cache.CACHE.IncrByFloat(ctx, "BALANCE#"+toAddress, value)

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timestamp)+8*60*60, 0).Format("2006-01-02")
		updateBalanceSnapshot(fromAddress, date, newFromBalance.Val())
		updateBalanceSnapshot(toAddress, date, newToBalance.Val())

		if fromAddress == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
			updateDailyBuyMetrics(toAddress, date, value)
		} else if toAddress == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
			updateDailySellMetrics(fromAddress, date, value)
		}
	}

	// 更新最后一次扫描的时间
	cache.CACHE.Set(ctx, "erc20_transfer_bsc_last_timestamp", newLastTimestamp, 0)

	return nil
}

// updateBalanceSnapshot 更新余额快照，便于每日归档
func updateBalanceSnapshot(address string, date string, value float64) {
	_, err := mysql.MYSQL.Exec(`INSERT INTO b_daily_onchain_trade_metrics (walletAddress, date, last_balance) VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE last_balance = VALUES(last_balance)`, address, date, value)
	if err != nil {
		log.Println("Failed to update balance snapshot for", address, "on", date)
		return
	}
	log.Println("Updated balance snapshot for", address, "on", date)
}

// updateDailyBuyMetrics 更新当天的够买数量和额度
func updateDailyBuyMetrics(address string, date string, value float64) {
	_, err := mysql.MYSQL.Exec(`
INSERT INTO b_daily_onchain_trade_metrics (walletAddress, date, buy_amount, buy_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE buy_amount = VALUES(buy_amount) + ?, buy_counts = VALUES(buy_counts) + ?
`, address, date, value, 1, value, 1)
	if err != nil {
		log.Println("Failed to update buy metrics for", address, "on", date)
		return
	}
	log.Println("Updated buy metrics for", address, "on", date)
}

// updateDailySellMetrics 更新当天的卖出数量和额度
func updateDailySellMetrics(address string, date string, value float64) {
	_, err := mysql.MYSQL.Exec(`
INSERT INTO b_daily_onchain_trade_metrics (walletAddress, date, sell_amount, sell_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE sell_amount = VALUES(sell_amount) + ?, sell_counts = VALUES(sell_counts) + ?
`, address, date, value, 1, value, 1)
	if err != nil {
		log.Println("Failed to update sell metrics for", address, "on", date)
		return
	}
	log.Println("Updated sell metrics for", address, "on", date)
}
