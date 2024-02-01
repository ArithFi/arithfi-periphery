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

func GenerateTxTag(from string, to string, amountETH *big.Float) string {
	fromNickname := "用户" + from[:7]
	toNickname := "用户" + from[:7]
	doWhat := "转账"
	howMuch := amountETH.Text('f', 2)
	fromIsDex := false
	toIsDex := false

	if from == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
		fromNickname = "在 PancakeSwap 上"
		doWhat = "买入"
		fromIsDex = true
	}

	if to == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
		toNickname = "在 PancakeSwap 上"
		doWhat = "卖出"
		toIsDex = true
	}

	if fromIsDex {
		return toNickname + " " + fromNickname + " " + doWhat + " " + howMuch
	} else if toIsDex {
		return fromNickname + " " + toNickname + " " + doWhat + " " + howMuch
	} else {
		return fromNickname + " " + doWhat + " " + howMuch + " 给 " + toNickname
	}
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
				log.Println("transfer_logs_analysis: Error loading topics")
				return
			}
			from, _ := topics[1].(string)
			from = "0x" + from[len(from)-40:]
			to, _ := topics[2].(string)
			to = "0x" + to[len(to)-40:]
			timestamp := new(big.Int)
			timestamp.SetString(strings.TrimPrefix(_log["timestamp"].(string), "0x"), 16)
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				log.Println("transfer_logs_analysis: Error loading location", err)
				return
			}
			localdate := time.Unix(timestamp.Int64(), 0).In(loc).Format("2006-01-02")
			localtime := time.Unix(timestamp.Int64(), 0).In(loc).Format("15:04:05")
			amountWei := new(big.Int)
			amountWei.SetString(strings.TrimPrefix(_log["data"].(string), "0x"), 16)
			amountEth := ConvertWeiToEth(amountWei)
			tag := GenerateTxTag(from, to, amountEth)

			abstract := bson.M{
				"from":   from,
				"to":     to,
				"amount": amountEth.String(),
				"tag":    tag,
			}
			aggregate := bson.M{
				"date":     localdate,
				"time":     localtime,
				"location": "Asia/Shanghai",
			}
			_, err = collection.UpdateOne(ctx, bson.M{"_id": _log["_id"]}, bson.M{"$set": bson.M{"abstract": abstract, "aggregate": aggregate}})
			if err != nil {
				return
			}
			log.Println("transfer_logs_analysis: success, block", _log["blockNumber"], ", date", localdate, ", time", localtime)
			log.Println(tag)
			fromBlock = _log["blockNumber"].(string)
		}
		log.Println("transfer_logs_analysis: Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
