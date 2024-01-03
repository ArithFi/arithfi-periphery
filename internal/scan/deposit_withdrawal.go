package scan

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

// DepositWithdrawal 扫描这个表的事件
func DepositWithdrawal(c echo.Context) error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "deposit_withdrawal_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	fmt.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT walletAddress, amount, timestamp, ordertype
FROM deposit_withdrawal 
WHERE timestamp > ?
AND status = 1
ORDER By timestamp 
LIMIT 100
`, lastTimestamp.Val())
	if err != nil {
		return err
	}
	defer query.Close()
	var newLastTimestamp int

	tx, err := mysql.MYSQL.Begin()

	for query.Next() {
		var walletAddress, ordertype string
		var timestamp int
		var amount float64
		err := query.Scan(&walletAddress, &amount, &timestamp, &ordertype)
		if err != nil {
			return err
		}
		fmt.Println("timestamp:", timestamp)
		newLastTimestamp = timestamp

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timestamp)+8*60*60, 0).Format("2006-01-02")

		if ordertype == "DEPOSIT" || ordertype == "WALLET_DEPOSIT" {
			handleDeposit(tx, walletAddress, amount, date)
		} else if ordertype == "WITHDRAW" {
			handleWithdraw(tx, walletAddress, amount, date)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	// 更新最后一次扫描的时间
	cache.CACHE.Set(ctx, "deposit_withdrawal_last_timestamp", newLastTimestamp, 0)

	return nil
}

// updateBalanceSnapshot 更新余额快照，便于每日归档
func handleDeposit(tx *sql.Tx, walletAddress string, amount float64, date string) {
	_, err := tx.Exec(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (walletAddress, date, deposit_amount, deposit_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE deposit_amount = VALUES(deposit_amount) + ?, deposit_counts = VALUES(deposit_counts) + ?`, walletAddress, date, amount, 1, amount, 1)
	if err != nil {
		log.Println("Failed to updates deposit snapshot for", walletAddress, "on", date)
		return
	}
	log.Println("Updated deposit snapshot for", walletAddress, "on", date)
}

// updateDailyBuyMetrics 更新当天的够买数量和额度
func handleWithdraw(tx *sql.Tx, walletAddress string, amount float64, date string) {
	_, err := tx.Exec(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (walletAddress, date, withdraw_amount, withdraw_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE withdraw_amount = VALUES(withdraw_amount) + ?, withdraw_counts = VALUES(withdraw_counts) + ?`, walletAddress, date, amount, 1, amount, 1)
	if err != nil {
		log.Println("Failed to update withdraw snapshot for", walletAddress, "on", date)
		return
	}
	log.Println("Updated withdraw snapshot for", walletAddress, "on", date)
}
