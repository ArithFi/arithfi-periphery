package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
	"strings"
	"time"
)

// This code can run continuously.
// The function is to retrieve transfer information from MongoDB, and analyze and expand the fields of the information.

// ConvertWeiToEth 将 wei 单位转换为 ETH 单位。
func ConvertWeiToEth(wei *big.Int) *big.Float {
	weiInEth := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(weiInEth, big.NewFloat(1e18))
	return ethValue
}

func main() {
	var fromBlock = "0"
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blocknumber", 1}}) // 按照blocknumber升序排序
	opts.SetLimit(200)

	for {
		collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
		cursor, err := collection.Find(ctx, bson.M{"blocknumber": bson.M{"$gte": fromBlock}}, opts)
		if err != nil {
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var log bson.M
			if err := cursor.Decode(&log); err != nil {
				continue
			}
			// 如果已经存在摘要，则先跳过该记录
			// 后续可以增加检查逻辑，验算摘要内容

			topics, ok := log["topics"].(bson.A)
			if !ok {
				fmt.Println(log["topics"])
				fmt.Println("无法获取topics字段或者topics字段不是切片类型")
				return
			}
			from, _ := topics[1].(string)
			to, _ := topics[2].(string)
			// 处理日期，北京时间，需要添加到摘要
			timeStamp := new(big.Int)
			timeStamp.SetString(strings.TrimPrefix(log["timestamp"].(string), "0x"), 16)
			loc, _ := time.LoadLocation("Asia/Shanghai")
			date := time.Unix(timeStamp.Int64(), 0).In(loc).Format("2006-01-02")
			// 转账金额，需要添加到摘要
			amountWei := new(big.Int)
			amountWei.SetString(strings.TrimPrefix(log["data"].(string), "0x"), 16)
			amountEth := ConvertWeiToEth(amountWei)
			// 更新 MongoDB 中的记录，增加摘要字段，方便后续分析
			// 创建一个新的对象，abstract: { from, to, date, amount}
			// 将摘要字段添加到 log 中
			abstract := bson.M{
				"from":   "0x" + from[len(from)-40:],
				"to":     "0x" + to[len(to)-40:],
				"amount": amountEth.String(),
			}
			aggregate := bson.M{
				"date":     date,
				"location": "Asia/Shanghai",
			}
			log["abstract"] = abstract
			log["aggregate"] = aggregate

			_, err := collection.UpdateOne(ctx, bson.M{"_id": log["_id"]}, bson.M{"$set": log})
			if err != nil {
				return
			}

			fmt.Println("更新记录成功", log["blocknumber"])
			fromBlock = log["blocknumber"].(string)
		}
		fmt.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10) // 每隔 10 秒获取一次记录
	}
}
