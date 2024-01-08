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
	opts.SetSort(bson.D{{"blocknumber", 1}})
	opts.SetLimit(2000)

	balancesMap := make(map[string]*big.Float)
	snapshotMap := make(map[string]map[string]*big.Float)
	lockMap := make(map[string]bool)
	totalTransfers := 0

	var snapshotCursorDate = "2023-09-25"

	for {
		fmt.Println("Start fetching data:", snapshotCursorDate)
		collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
		cursor, err := collection.Find(ctx, bson.M{"aggregate.date": bson.M{"$gte": snapshotCursorDate}}, opts)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var log bson.M
			if err := cursor.Decode(&log); err != nil {
				fmt.Println(err)
				continue
			}
			if lockMap[log["_id"].(string)] {
				continue
			}
			lockMap[log["_id"].(string)] = true
			aggregate, ok := log["aggregate"].(bson.M)
			if !ok {
				fmt.Println("Unable to retrieve the aggregate field or the aggregate field is not of slice type.")
				return
			}
			date := aggregate["date"].(string)
			totalTransfers++
			abstract, ok := log["abstract"].(bson.M)
			if !ok {
				fmt.Println("Unable to retrieve the abstract field or the abstract field is not of slice type.")
				return
			}
			from := abstract["from"].(string)
			to := abstract["to"].(string)
			amount, ok := new(big.Float).SetString(abstract["amount"].(string))
			if !ok {
				fmt.Println(abstract["amount"])
				fmt.Println("Unable to retrieve the amount field or the amount field is not of string type.")
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

			var snapshotArray []bson.M
			for address, balance := range snapshotMap[date] {
				if balance.Cmp(big.NewFloat(0)) > 0 {
					snapshotArray = append(snapshotArray, bson.M{"address": address, "quantity": balance.String(), "percentage": new(big.Float).Quo(balance, totalSupply).String()})
				}
			}

			if date > snapshotCursorDate {
				var _abstract bson.M
				_abstract = make(bson.M)
				_abstract["holders"] = len(snapshotArray)
				_abstract["total_transfers"] = totalTransfers
				collection := mongo.MONGODB.Database("chain-bsc").Collection("tokenholder-snapshot")
				_, err := collection.UpdateOne(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"abstract": _abstract, "holders": snapshotArray}}, options.Update().SetUpsert(true))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Snapshot updated:", date, "holders:", len(snapshotArray), "total_transfers:", totalTransfers)
			}
			snapshotCursorDate = date
		}
		fmt.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
