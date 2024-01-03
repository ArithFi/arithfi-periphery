package scan

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"log"
	"strconv"
	"time"
)

// DepositWithdrawal 扫描这个表的事件
func DepositWithdrawal() error {
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "deposit_withdrawal_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	log.Println(lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT walletAddress, amount, timestamp, ordertype
FROM deposit_withdrawal 
WHERE timestamp > ?
AND status = 1
ORDER By timestamp 
LIMIT 200
`, lastTimestamp.Val())
	if err != nil {
		return err
	}
	defer query.Close()
	var newLastTimestamp int

	tx, err := mysql.MYSQL.Begin()
	handleDepositStmt, err := tx.Prepare(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (walletAddress, date, deposit_amount, deposit_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE deposit_amount = VALUES(deposit_amount) + deposit_amount, deposit_counts = VALUES(deposit_counts) + deposit_counts`)
	handleWithdrawStmt, err := tx.Prepare(`INSERT INTO b_daily_offchain_deposit_withdraw_metrics (walletAddress, date, withdraw_amount, withdraw_counts) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE withdraw_amount = VALUES(withdraw_amount) + withdraw_amount, withdraw_counts = VALUES(withdraw_counts) + withdraw_counts`)
	for query.Next() {
		var walletAddress, ordertype string
		var timestamp int
		var amount float64
		err := query.Scan(&walletAddress, &amount, &timestamp, &ordertype)
		if err != nil {
			return err
		}
		log.Println("timestamp:", timestamp)
		newLastTimestamp = timestamp

		// 获取时间戳，需要处理成+8的北京时间,获取北京的时间的日期字符串
		date := time.Unix(int64(timestamp)+8*60*60, 0).Format("2006-01-02")

		if ordertype == "DEPOSIT" || ordertype == "WALLET_DEPOSIT" {
			_, err = handleDepositStmt.Exec(walletAddress, date, amount, 1)
			if err != nil {
				return err
			}
		} else if ordertype == "WITHDRAW" {
			_, err = handleWithdrawStmt.Exec(walletAddress, date, amount, 1)
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
		cache.CACHE.Set(ctx, "deposit_withdrawal_last_timestamp", newLastTimestamp, 0)
	}

	return nil
}
