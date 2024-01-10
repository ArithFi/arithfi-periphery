package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"strings"
	"time"
)

// This code can run continuously.
// The function is to retrieve transfer information from MongoDB, and analyze and expand the fields of the information.

// ConvertWeiToEth is used to convert wei to eth.
func ConvertWeiToEth(wei *big.Int) *big.Float {
	weiInEth := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(weiInEth, big.NewFloat(1e18))
	return ethValue
}

func main() {
	var fromBlock = "0"
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blockNumber", 1}}) // 按照blockNumber升序排序
	opts.SetLimit(1000)

	collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")

	for {
		cursor, err := collection.Find(ctx, bson.M{"blockNumber": bson.M{"$gte": fromBlock}}, opts)
		if err != nil {
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var _log bson.M
			if err := cursor.Decode(&_log); err != nil {
				continue
			}
			topics, ok := _log["topics"].(bson.A)
			if !ok {
				log.Println("无法获取topics字段或者topics字段不是切片类型")
				return
			}
			from, _ := topics[1].(string)
			to, _ := topics[2].(string)
			timestamp := new(big.Int)
			timestamp.SetString(strings.TrimPrefix(_log["timestamp"].(string), "0x"), 16)
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				log.Println("Error loading location:", err)
				return
			}
			date := time.Unix(timestamp.Int64(), 0).In(loc).Format("2006-01-02")
			amountWei := new(big.Int)
			amountWei.SetString(strings.TrimPrefix(_log["data"].(string), "0x"), 16)
			amountEth := ConvertWeiToEth(amountWei)
			abstract := bson.M{
				"from":   "0x" + from[len(from)-40:],
				"to":     "0x" + to[len(to)-40:],
				"amount": amountEth.String(),
			}
			aggregate := bson.M{
				"date":     date,
				"location": "Asia/Shanghai",
			}
			_, err = collection.UpdateOne(ctx, bson.M{"_id": _log["_id"]}, bson.M{"$set": bson.M{"abstract": abstract, "aggregate": aggregate}})
			if err != nil {
				return
			}

			log.Println("Update transfer_logs success, block:", _log["blockNumber"], ", date:", date)
			fromBlock = _log["blockNumber"].(string)
		}
		log.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
