package scan

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/labstack/echo/v4"
)

func ERC20TransferBSC(c echo.Context) error {
	// 遍历数据库100条数据
	// 根据交易来更新用户的余额

	// 根据交易来更新用户的日报告

	// 获取当前的扫描时间戳，默认是0
	var ctx = context.Background()
	lastTimestamp := cache.CACHE.Get(ctx, "erc20_transfer_bsc_last_timestamp")
	if lastTimestamp == nil {
		lastTimestamp.SetVal("0")
	}

	fmt.Println("erc20_transfer_bsc last timestamp:", lastTimestamp)

	query, err := mysql.MYSQL.Query(`SELECT from_address, to_address, timestamp, value
FROM erc20_transfer_bsc 
WHERE timestamp > ? 
ORDER By timestamp 
LIMIT 1
`, lastTimestamp.String())
	if err != nil {
		return err
	}
	defer query.Close()
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
		fmt.Println("value:", value)
	}

	return nil
}
