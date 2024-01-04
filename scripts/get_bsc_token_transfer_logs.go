package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/internal/bscscan"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"strings"
	"time"
)

func main() {
	var fromBlock = "0"
	const toBlock = "latest"

	for {
		fmt.Printf("fromBlock: %s, toBlock: %s\n", fromBlock, toBlock)
		logs, err := bscscan.GetLogs(fromBlock, toBlock)
		if err != nil {
			log.Fatalf("Error getting logs: %v", err)
		}
		var documents []interface{}
		for _, _log := range logs {
			documents = append(documents, _log)
			fromBlock = _log.BlockNumber
		}

		if len(documents) > 0 {
			fromBlockBigInt, ok := new(big.Int).SetString(strings.TrimPrefix(fromBlock, "0x"), 16)
			if !ok {
				log.Fatalf("Failed to parse fromBlock: %s", fromBlock)
			}
			fromBlockBigInt = fromBlockBigInt.Add(fromBlockBigInt, big.NewInt(1))
			fromBlock = fromBlockBigInt.String()

			ctx := context.TODO()
			collection := mongo.MONGODB.Database("chain-bsc").Collection("transfer-logs")
			_options := options.InsertMany().SetOrdered(false)
			_, err = collection.InsertMany(ctx, documents, _options)
			if err != nil {
				fmt.Println("Error inserting some documents")
			}
			fmt.Println("Logs inserted successfully")
		} else {
			fmt.Println("No logs to insert")
		}

		// add a delay or condition to exit the loop if needed
		time.Sleep(time.Second * 10)
	}
}
