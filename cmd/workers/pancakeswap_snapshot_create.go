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
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blockNumber", 1}})
	opts.SetLimit(2000)

	// date => totalBuyVolume
	totalBuyVolumeMap := make(map[string]*big.Float)
	// date => totalSellVolume
	totalSellVolumeMap := make(map[string]*big.Float)
	// date => totalBuyTransfers
	totalBuyTxsMap := make(map[string]int)
	// date => totalSellTransfers
	totalSellTxsMap := make(map[string]int)

	lockMap := make(map[string]bool)

	// date => address => {totalBuyVolume, totalSellVolume, totalBuyTxs, totalSellTxs}
	snapshotMap := make(map[string]map[string]map[string]*big.Float)

	var snapshotCursorDate = "2023-09-26"
	transferLogsCollection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
	pancakeSwapSnapshotCollection := mongo.MONGODB.Database("chain-bsc").Collection("pancakeswap-snapshot")

	for {
		log.Println("Start fetching data:", snapshotCursorDate)
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
				log.Println("Unable to retrieve the _id field or the _id field is not of ObjectID type.")
				continue
			}
			aggregate, ok := _log["aggregate"].(bson.M)
			if !ok {
				log.Println("Unable to retrieve the aggregate field or the aggregate field is not of slice type.")
				return
			}
			date := aggregate["date"].(string)
			abstract, ok := _log["abstract"].(bson.M)
			if !ok {
				log.Println("Unable to retrieve the abstract field or the abstract field is not of slice type.")
				return
			}
			from := abstract["from"].(string)
			to := abstract["to"].(string)
			amount, ok := new(big.Float).SetString(abstract["amount"].(string))
			if !ok {
				log.Println("Unable to retrieve the amount field or the amount field is not of string type.")
				return
			}
			if from == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" { // This is the address of the PancakeSwap contract, Buy
				if totalBuyVolumeMap[date] == nil {
					totalBuyVolumeMap[date] = new(big.Float)
				}
				totalBuyVolumeMap[date].Add(totalBuyVolumeMap[date], amount)
				totalBuyTxsMap[date]++
				if snapshotMap[date] == nil {
					snapshotMap[date] = make(map[string]map[string]*big.Float)
				}
				if snapshotMap[date][to] == nil {
					snapshotMap[date][to] = make(map[string]*big.Float)
				}
				if snapshotMap[date][to]["totalSellVolume"] == nil {
					snapshotMap[date][to]["totalSellVolume"] = new(big.Float)
				}
				if snapshotMap[date][to]["totalBuyVolume"] == nil {
					snapshotMap[date][to]["totalBuyVolume"] = new(big.Float)
				}
				snapshotMap[date][to]["totalBuyVolume"].Add(snapshotMap[date][to]["totalBuyVolume"], amount)
				if snapshotMap[date][to]["totalBuyTxs"] == nil {
					snapshotMap[date][to]["totalBuyTxs"] = new(big.Float)
				}
				snapshotMap[date][to]["totalBuyTxs"].Add(snapshotMap[date][to]["totalBuyTxs"], big.NewFloat(1))
			} else if to == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" { // This is the address of the PancakeSwap contract, Sell
				if totalSellVolumeMap[date] == nil {
					totalSellVolumeMap[date] = new(big.Float)
				}
				totalSellVolumeMap[date].Add(totalSellVolumeMap[date], amount)
				totalSellTxsMap[date]++
				if snapshotMap[date] == nil {
					snapshotMap[date] = make(map[string]map[string]*big.Float)
				}
				if snapshotMap[date][from] == nil {
					snapshotMap[date][from] = make(map[string]*big.Float)
				}
				if snapshotMap[date][from]["totalSellVolume"] == nil {
					snapshotMap[date][from]["totalSellVolume"] = new(big.Float)
				}
				snapshotMap[date][from]["totalSellVolume"].Add(snapshotMap[date][from]["totalSellVolume"], amount)
				if snapshotMap[date][from]["totalSellTxs"] == nil {
					snapshotMap[date][from]["totalSellTxs"] = new(big.Float)
				}
				snapshotMap[date][from]["totalSellTxs"].Add(snapshotMap[date][from]["totalSellTxs"], big.NewFloat(1))
			}

			if from == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" || to == "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38" {
				tradersArray := make([]bson.M, 0)
				for address, metrics := range snapshotMap[date] {
					if metrics["totalSellVolume"] == nil {
						metrics["totalSellVolume"] = new(big.Float)
					}
					if metrics["totalBuyVolume"] == nil {
						metrics["totalBuyVolume"] = new(big.Float)
					}
					if metrics["totalSellTxs"] == nil {
						metrics["totalSellTxs"] = new(big.Float)
					}
					if metrics["totalBuyTxs"] == nil {
						metrics["totalBuyTxs"] = new(big.Float)
					}
					totalSellVolume, _ := metrics["totalSellVolume"].Float64()
					totalBuyVolume, _ := metrics["totalBuyVolume"].Float64()
					totalSellTxs, _ := metrics["totalSellTxs"].Int64()
					totalBuyTxs, _ := metrics["totalBuyTxs"].Int64()
					tradersArray = append(tradersArray, bson.M{
						"address":               address,
						"totalSellVolume":       totalSellVolume,
						"totalBuyVolume":        totalBuyVolume,
						"totalVolume":           totalSellVolume + totalBuyVolume,
						"totalSellTransactions": totalSellTxs,
						"totalBuyTransactions":  totalBuyTxs,
						"totalTransactions":     totalBuyTxs + totalSellTxs,
					})
				}
				var _abstract bson.M
				_abstract = make(bson.M)
				_abstract["traders"] = len(tradersArray)
				if totalBuyVolumeMap[date] == nil {
					totalBuyVolumeMap[date] = new(big.Float)
				}
				if totalSellVolumeMap[date] == nil {
					totalSellVolumeMap[date] = new(big.Float)
				}
				_abstract["totalBuyTransactions"] = totalBuyTxsMap[date]
				_abstract["totalSellTransactions"] = totalSellTxsMap[date]
				_abstract["totalTransactions"] = totalBuyTxsMap[date] + totalSellTxsMap[date]
				_abstract["totalBuyVolume"], _ = totalBuyVolumeMap[date].Float64()
				_abstract["totalSellVolume"], _ = totalSellVolumeMap[date].Float64()
				_abstract["totalVolume"], _ = new(big.Float).Add(totalSellVolumeMap[date], totalBuyVolumeMap[date]).Float64()
				_, err := pancakeSwapSnapshotCollection.UpdateOne(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"abstract": _abstract, "traders": tradersArray}}, options.Update().SetUpsert(true))
				if err != nil {
					log.Println(err)
				}
				log.Println("Snapshot updated:", date, "from:", from, "to:", to)
			} else {
				log.Println("None pancakeswap transactions")
			}
			snapshotCursorDate = date
		}
		log.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
