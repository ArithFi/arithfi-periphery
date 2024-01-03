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
	lastId := cache.CACHE.Get(ctx, "f_future_trading_last_id")
	if lastId == nil {
		lastId.SetVal("0")
	}

	log.Println(lastId)

	query, err := mysql.MYSQL.Query(`SELECT _id, timeStamp, product, positionIndex, leverage, orderType, mode, direction, margin, volume, sellValue, walletAddress, kolAddress
FROM f_future_trading 
WHERE _id > ? 
ORDER By _id 
LIMIT 200
`, lastId.Val())
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
		var timeStamp string
		var id int
		var leverage int64
		var orderType string
		var mode string
		var direction string
		var margin float64
		var volume float64
		var sellValue float64
		var walletAddress string
		var kolAddress string

		err := query.Scan(&id, &timeStamp, &product, &positionIndex, &leverage, &orderType, &mode, &direction, &margin, &volume, &sellValue, &walletAddress, &kolAddress)
		if err != nil {
			return err
		}
		newLastId = id
		log.Println(id, product, positionIndex, leverage, orderType, mode, direction, margin, volume, sellValue, walletAddress, kolAddress)

		loc, _ := time.LoadLocation("Local")
		date, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStamp, loc)

		// 将时间对象转换为北京时间 (+8)
		beijingLoc, _ := time.LoadLocation("Asia/Shanghai")
		beijingDate := date.In(beijingLoc)
		beijingDateStr := beijingDate.Format("2006-01-02")
		log.Println("date:", beijingDateStr)
		if orderType == "MARKET_ORDER_FEE" || orderType == "LIMIT_ORDER_FEE" {
			_, err = handleNewOrderStmt.Exec(beijingDateStr, walletAddress, mode, kolAddress, 1, volume)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Fatalf("insert error: %v, unable to rollback: %v", err, rbErr)
				}
				return err
			}
		} else if orderType == "MARKET_CLOSE_FEE" || orderType == "TP_ORDER_FEE" || orderType == "SL_ORDER_FEE" || orderType == "MARKET_LIQUIDATION" {
			_, err = handleBurnStmt.Exec(beijingDateStr, walletAddress, mode, kolAddress, sellValue-margin)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					log.Fatalf("insert error: %v, unable to rollback: %v", err, rbErr)
				}
				return err
			}
		} else {
			break
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
