package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"time"
)

func main() {
	totalSupply := new(big.Float)
	totalSupply.SetString("300000000")

	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blockNumber", 1}})
	opts.SetLimit(2000)

	balancesMap := make(map[string]*big.Float)
	snapshotMap := make(map[string]map[string]*big.Float)
	lockMap := make(map[string]bool)
	totalTransfers := 0

	var snapshotCursorDate = "2023-09-25"

	transferLogsCollection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
	tokenHolderSnapshotCollection := mongo.MONGODB.Database("chain-bsc").Collection("tokenholder-snapshot")

	for {
		log.Println("tokenholder_snapshot_create: Start fetching data", snapshotCursorDate)
		cursor, err := transferLogsCollection.Find(ctx, bson.M{"aggregate.date": bson.M{"$gte": snapshotCursorDate}}, opts)
		if err != nil {
			log.Println(err)
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var _log bson.M
			if err := cursor.Decode(&_log); err != nil {
				log.Println(err)
				continue
			}
			if id, ok := _log["_id"].(primitive.ObjectID); ok {
				idHex := id.Hex()
				if lockMap[idHex] {
					continue
				} else {
					lockMap[idHex] = true
				}
			} else {
				log.Println("tokenholder_snapshot_create: Unable to retrieve the _id field or the _id field is not of ObjectID type.")
				continue
			}
			aggregate, ok := _log["aggregate"].(bson.M)
			if !ok {
				log.Println("tokenholder_snapshot_create: Unable to retrieve the aggregate field or the aggregate field is not of slice type.")
				return
			}
			date := aggregate["date"].(string)
			totalTransfers++
			abstract, ok := _log["abstract"].(bson.M)
			if !ok {
				log.Println("tokenholder_snapshot_create: Unable to retrieve the abstract field or the abstract field is not of slice type.")
				return
			}
			from := abstract["from"].(string)
			to := abstract["to"].(string)
			amount, ok := new(big.Float).SetString(abstract["amount"].(string))
			if !ok {
				log.Println("tokenholder_snapshot_create: Unable to retrieve the amount field or the amount field is not of string type.")
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

			var _abstract bson.M
			_abstract = make(bson.M)
			_abstract["holders"] = len(snapshotArray)
			_abstract["transfers"] = bson.M{
				"total": totalTransfers,
			}
			_, err := tokenHolderSnapshotCollection.UpdateOne(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"abstract": _abstract, "holders": snapshotArray}}, options.Update().SetUpsert(true))
			if err != nil {
				log.Println(err)
			}
			log.Println("tokenholder_snapshot_create: success", date, "holders", len(snapshotArray), "totalTransfers", totalTransfers)

			if date > snapshotCursorDate {
				delete(snapshotMap, snapshotCursorDate)
				log.Println("tokenholder_snapshot_create: delete snapshotMap of", snapshotCursorDate)
			}
			snapshotCursorDate = date
		}
		log.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
