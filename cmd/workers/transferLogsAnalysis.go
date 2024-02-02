package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/arithfi/arithfi-periphery/internal/bscscan"
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
	db := mysql.ArithFiDB

	UserTagMap := bscscan.UserMap{
		"0xdccbdbaee4d9d6639242f18f4eb08f4edad1a331": "ArithFi: System",
		"0x7c4fb3E5ba0a5D80658889715b307e66916f29b2": "ArithFi: Deployer",
		"0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38": "PancakeSwap",
		"0x0000000000000000000000000000000000000000": "NULL",
		"0xe26d976910D688083c8F9eCcB25e42345E5b95a0": "ArithFi: BSC-ETH-Bridge",
	}

	query, err := db.Query(`SELECT walletAddress, type, tgName, country FROM f_kol_info`)
	if err != nil {
		return
	}
	for query.Next() {
		var walletAddress, typeStr, tgName, country string
		if err := query.Scan(&walletAddress, &typeStr, &tgName, &country); err != nil {
			continue
		}
		walletAddress = strings.ToLower(walletAddress)
		if tgName == "" {
			tgName = walletAddress[0:6]
		}
		if UserTagMap[walletAddress] != "" {
			continue
		}
		UserTagMap[walletAddress] = typeStr + "-" + tgName + "-" + country
	}

	query, err = db.Query(`SELECT walletAddress FROM f_user_assets`)
	if err != nil {
		return
	}
	for query.Next() {
		var walletAddress string
		if err := query.Scan(&walletAddress); err != nil {
			continue
		}
		walletAddress = strings.ToLower(walletAddress)
		if UserTagMap[walletAddress] != "" {
			continue
		}
		UserTagMap[walletAddress] = "Trader"
	}

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
			tag := bscscan.GenerateTxTag(from, to, amountEth, UserTagMap)

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
