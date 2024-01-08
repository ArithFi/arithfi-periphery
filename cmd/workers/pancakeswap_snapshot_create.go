package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
	"time"
)

func main() {
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"blocknumber", 1}})
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
			if id, ok := log["_id"].(primitive.ObjectID); ok {
				idHex := id.Hex()
				if lockMap[idHex] {
					continue
				} else {
					lockMap[idHex] = true
				}
			} else {
				fmt.Println("Unable to retrieve the _id field or the _id field is not of ObjectID type.")
				continue
			}
			aggregate, ok := log["aggregate"].(bson.M)
			if !ok {
				fmt.Println("Unable to retrieve the aggregate field or the aggregate field is not of slice type.")
				return
			}
			date := aggregate["date"].(string)
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
					totalSellVolume := metrics["totalSellVolume"]
					totalBuyVolume := metrics["totalBuyVolume"]
					totalSellTxs, _ := metrics["totalSellTxs"].Int64()
					totalBuyTxs, _ := metrics["totalBuyTxs"].Int64()
					tradersArray = append(tradersArray, bson.M{
						"address":                 address,
						"total_sell_volume":       totalSellVolume.String(),
						"total_buy_volume":        totalBuyVolume.String(),
						"total_volume":            new(big.Float).Add(totalSellVolume, totalBuyVolume).String(),
						"total_sell_transactions": totalSellTxs,
						"total_buy_transactions":  totalBuyTxs,
						"total_transactions":      totalBuyTxs + totalSellTxs,
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
				_abstract["total_buy_transactions"] = totalBuyTxsMap[date]
				_abstract["total_sell_transactions"] = totalSellTxsMap[date]
				_abstract["total_transactions"] = totalBuyTxsMap[date] + totalSellTxsMap[date]
				_abstract["total_buy_volume"] = totalBuyVolumeMap[date].String()
				_abstract["total_sell_volume"] = totalSellVolumeMap[date].String()
				abstract["total_volume"] = new(big.Float).Add(totalSellVolumeMap[date], totalBuyVolumeMap[date]).String()
				collection := mongo.MONGODB.Database("chain-bsc").Collection("pancakeswap-snapshot")
				_, err := collection.UpdateOne(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"abstract": _abstract, "traders": tradersArray}}, options.Update().SetUpsert(true))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Snapshot updated:", date, "from:", from, "to:", to)
			} else {
				fmt.Println("None pancakeswap transactions")
			}
			snapshotCursorDate = date
		}
		fmt.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}