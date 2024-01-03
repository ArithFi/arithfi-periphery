package scan

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"log"
	"strconv"
	"time"
)

// ERC20TransferBSC 扫描这个表的事件
func ERC20TransferBSC() error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "erc20_transfer_bsc_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	log.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT from_address, to_address, timestamp, value
FROM erc20_transfer_bsc 
WHERE timestamp > ? 
ORDER By timestamp
LIMIT 200
`, lastTimestamp.Val())
	if err != nil {
		return err
	}
	defer query.Close()
	newLastTimestamp := 0

	tx, err := mysql.MYSQL.Begin()
	updateDailyBuyMetricsStmt, err := tx.Prepare(`INSERT INTO b_daily_onchain_trade_metrics (walletAddress, date, buy_amount, buy_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE buy_amount = VALUES(buy_amount) + buy_amount, buy_counts = VALUES(buy_counts) + buy_counts`)
	updateDailySellMetricsStmt, err := tx.Prepare(`INSERT INTO b_daily_onchain_trade_metrics (walletAddress, date, sell_amount, sell_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE sell_amount = VALUES(sell_amount) + sell_amount, sell_counts = VALUES(sell_counts) + sell_counts`)
	for query.Next() {
		var fromAddress, toAddress string
		var timestamp int
		var value float64
		err := query.Scan(&fromAddress, &toAddress, &timestamp, &value)
		if err != nil {
			return err
		}
		log.Println("fromAddress:", fromAddress)
		log.Println("toAddress:", toAddress)
		log.Println("timestamp:", timestamp)
		newLastTimestamp = timestamp
		log.Println("value:", value)

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timestamp)+8*60*60, 0).Format("2006-01-02")

		if fromAddress == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
			_, err = updateDailyBuyMetricsStmt.Exec(toAddress, date, value, 1)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Fatalf("insert error: %v, unable to rollback: %v", err, rbErr)
				}
				return err
			}
		} else if toAddress == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
			_, err = updateDailySellMetricsStmt.Exec(fromAddress, date, value, 1)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Fatalf("insert error: %v, unable to rollback: %v", err, rbErr)
				}
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	lastTimestampNumber, err := strconv.Atoi(lastTimestamp.Val())
	if err != nil {
		lastTimestampNumber = 0
	}
	if newLastTimestamp > lastTimestampNumber {
		cache.CACHE.Set(ctx, "erc20_transfer_bsc_last_timestamp", newLastTimestamp, 0)
	}

	return nil
}
