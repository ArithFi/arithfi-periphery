package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
	"time"
)

func main() {
	totalSupply := new(big.Float)
	totalSupply.SetString("300000000")

	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blocknumber", 1}}) // 按照blocknumber升序排序
	opts.SetLimit(2000)

	// 维护一个map，用于存储每个地址的总额
	balancesMap := make(map[string]*big.Float)
	snapshotMap := make(map[string]map[string]*big.Float)
	totalTransfersMap := make(map[string]int)

	// 每次更新 snapshot 后更新 snapshotCursorDate
	var snapshotCursorDate = "2023-09-25"

	for {
		fmt.Println("开始获取数据:", snapshotCursorDate)
		collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
		cursor, err := collection.Find(ctx, bson.M{"aggregate.date": bson.M{"$gte": snapshotCursorDate}}, opts)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer cursor.Close(ctx)
		fmt.Println("准备遍历数据")
		for cursor.Next(ctx) {
			var log bson.M
			if err := cursor.Decode(&log); err != nil {
				fmt.Println(err)
				continue
			}
			aggregate, ok := log["aggregate"].(bson.M)
			if !ok {
				fmt.Println(log["aggregate"])
				fmt.Println("无法获取aggregate字段或者aggregate字段不是切片类型")
				return
			}
			date := aggregate["date"].(string)
			totalTransfersMap[date] = totalTransfersMap[date] + 1
			fmt.Println("新总转账次数", date, ":", totalTransfersMap[date])
			abstract, ok := log["abstract"].(bson.M)
			if !ok {
				fmt.Println(log["abstract"])
				fmt.Println("无法获取abstract字段或者abstract字段不是切片类型")
				return
			}
			from := abstract["from"].(string)
			to := abstract["to"].(string)
			// 使用 *big.Float 处理大数
			amount, ok := new(big.Float).SetString(abstract["amount"].(string))
			if !ok {
				fmt.Println(abstract["amount"])
				fmt.Println("无法获取amount字段或者amount字段不是字符串类型")
				return
			}
			if balancesMap[from] == nil {
				balancesMap[from] = new(big.Float)
			}
			if from != "0x0000000000000000000000000000000000000000" {
				balancesMap[from].Sub(balancesMap[from], amount)
			}
			if balancesMap[to] == nil {
				balancesMap[to] = new(big.Float)
			}
			balancesMap[to].Add(balancesMap[to], amount)
			snapshotMap[date] = balancesMap

			if date > snapshotCursorDate {
				var snapshotArray []bson.M
				for address, balance := range snapshotMap[snapshotCursorDate] {
					// balance > 0, append to snapshot array
					if balance.Cmp(big.NewFloat(0)) > 0 {
						snapshotArray = append(snapshotArray, bson.M{"address": address, "quantity": balance.String(), "percentage": new(big.Float).Quo(balance, totalSupply).String()})
					}
				}
				var abstract bson.M
				abstract = make(bson.M)
				abstract["holders"] = len(snapshotArray)
				abstract["total_transfers"] = totalTransfersMap[snapshotCursorDate]
				collection := mongo.MONGODB.Database("chain-bsc").Collection("tokenholder-snapshot")
				_, err := collection.UpdateOne(ctx, bson.M{"date": snapshotCursorDate}, bson.M{"$set": bson.M{"abstract": abstract, "holders": snapshotArray}}, options.Update().SetUpsert(true))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Snapshot updated:", snapshotCursorDate)
				// update snapshotCursorDate
				snapshotCursorDate = date
			}
		}
		fmt.Println("Sleep 5 seconds")
		time.Sleep(time.Second * 5) // 每隔 5 秒获取一次记录
	}
}
