package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/internal/bscscan"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"strings"
	"time"
)

// This code can execute continuously in a loop, triggering every 10 seconds.
// The purpose of the code is to retrieve transfer records of ATF tokens from the Logs data source on BscScan, all of
// which are successful transactions. The logs will be continuously written into MongoDB, where there is a unique
// constraint: a combination of transactionHash and logIndex forms a unique index.

func main() {
	var fromBlock = "0"
	const toBlock = "latest"

	ctx := context.TODO()
	collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")

	for {
		log.Printf("transfer_logs_retrieve: fromBlock %s, toBlock%s\n", fromBlock, toBlock)
		logs, err := bscscan.GetLogs(fromBlock, toBlock)
		if err != nil {
			log.Fatalf("transfer_logs_retrieve: Error getting logs %v", err)
		}
		var documents []interface{}
		for _, _log := range logs {
			documents = append(documents, _log)
			fromBlock = _log.BlockNumber
		}

		if len(documents) > 0 {
			fromBlockBigInt, ok := new(big.Int).SetString(strings.TrimPrefix(fromBlock, "0x"), 16)
			if !ok {
				log.Fatalf("transfer_logs_retrieve: Failed to parse fromBlock %s", fromBlock)
			}
			fromBlockBigInt = fromBlockBigInt.Add(fromBlockBigInt, big.NewInt(1))
			fromBlock = fromBlockBigInt.String()

			_options := options.InsertMany().SetOrdered(false)
			_, err = collection.InsertMany(ctx, documents, _options)
			if err != nil {
				log.Println("Error inserting some documents")
			}
			log.Println("transfer_logs_retrieve success")
		} else {
			log.Println("transfer_logs_retrieve: No logs to insert")
		}

		log.Println("transfer_logs_retrieve: Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
