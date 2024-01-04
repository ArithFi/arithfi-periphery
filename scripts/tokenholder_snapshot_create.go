package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

func main() {
	var fromBlock = "0"
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blocknumber", 1}}) // 按照blocknumber升序排序
	opts.SetLimit(200)

	// 维护一个map，用于存储每个地址的总额
	balancesMap := make(map[string]float64)
	snapshotMap := make(map[string]map[string]float64)

	// 每次更新 snapshot 后更新 snapshotCursorDate
	var snapshotCursorDate = ""

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
			aggregate, ok := log["aggregate"].(bson.M)
			if !ok {
				fmt.Println(log["aggregate"])
				fmt.Println("无法获取aggregate字段或者aggregate字段不是切片类型")
				return
			}
			date := aggregate["date"].(string)
			abstract, ok := log["abstract"].(bson.M)
			if !ok {
				fmt.Println(log["abstract"])
				fmt.Println("无法获取abstract字段或者abstract字段不是切片类型")
				return
			}
			from := abstract["from"].(string)
			to := abstract["to"].(string)
			amount, err := strconv.ParseFloat(abstract["amount"].(string), 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			if from != "0x0000000000000000000000000000000000000000" {
				balancesMap[from] -= amount
			}
			balancesMap[to] += amount
			snapshotMap[date] = balancesMap

			if date > snapshotCursorDate {
				// 准备更新snapshot
				fmt.Println("准备更新snapshot", date)
				fmt.Println(snapshotMap[date])
				snapshotCursorDate = date
			}
		}
		fmt.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10) // 每隔 10 秒获取一次记录
	}
}
