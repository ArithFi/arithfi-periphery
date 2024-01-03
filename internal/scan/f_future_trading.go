package scan

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"log"
	"time"
)

// FFutureTrading 扫描这个表
func FFutureTrading() error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "f_future_trading_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	log.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT _id, product, positionIndex, leverage, orderType, mode, direction, margin, volume, sellValue, walletAddress, kolAddress, availableBanlance, copyAccountBalance
FROM f_future_trading 
WHERE _id > ? 
ORDER By _id 
LIMIT 200
`, lastTimestamp.Val())
	if err != nil {
		return err
	}
	defer query.Close()
	var newLastId int
	tx, err := mysql.MYSQL.Begin()
	handleNewOrderStmt, err := tx.Prepare(`INSERT INTO b_daily_offchain_futures_metrics (date, walletAddress, mode, kolAddress, new_position_counts, new_position_size)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE new_position_counts = new_position_counts + VALUES(new_position_counts), new_position_size = new_position_size + VALUES(new_position_size)`)
	handleBurnStmt, err := tx.Prepare(`INSERT INTO b_daily_offchain_futures_metrics (date, walletAddress, mode, kolAddress, net_burn_amount)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE net_burn_amount = VALUES(net_burn_amount) + net_burn_amount`)
	for query.Next() {
		var product string
		var positionIndex int64
		var timeStamp, id int
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

		err := query.Scan(&id, &product, &positionIndex, &leverage, &orderType, &mode, &direction, &margin, &volume, &sellValue, &walletAddress, &kolAddress, &availableBanlance, &copyAccountBalance)
		if err != nil {
			return err
		}
		newLastId = id

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timeStamp)+8*60*60, 0).Format("2006-01-02")
		log.Println("date:", date)
		if orderType == "MARKET_ORDER_FEE" || orderType == "LIMIT_ORDER_FEE" {
			_, err = handleNewOrderStmt.Exec(date, walletAddress, mode, kolAddress, 1, volume, volume)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Fatalf("insert error: %v, unable to rollback: %v", err, rbErr)
				}
				return err
			}
		} else if orderType == "MARKET_CLOSE_FEE" || orderType == "TP_ORDER_FEE" || orderType == "SL_ORDER_FEE" || orderType == "MARKET_LIQUIDATION" {
			_, err = handleBurnStmt.Exec(date, walletAddress, mode, kolAddress, sellValue-margin)
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
	if newLastId > 0 {
		cache.CACHE.Set(ctx, "f_future_trading_last_id", newLastId, 0)
	}
	return nil
}
