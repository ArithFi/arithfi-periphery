package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// This code can run continuously.
// The function is to retrieve transfer information from MongoDB, and analyze and expand the fields of the information.

func main() {
	var fromPositionIndex = int64(-1)
	ctx := context.TODO()

	opts := options.Find()
	opts.SetSort(bson.D{{"positionIndex", 1}}) // 按照blockNumber升序排序
	opts.SetLimit(1000)

	collection := mongo.MONGODB.Database("off-chain").Collection("futures-positions")

	for {
		cursor, err := collection.Find(ctx, bson.M{"positionIndex": bson.M{"$gt": fromPositionIndex}}, opts)
		if err != nil {
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var position bson.M
			if err := cursor.Decode(&position); err != nil {
				continue
			}
			timestamp := position["timestamp"].(int64)
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				log.Println("Error loading location:", err)
				return
			}
			date := time.Unix(timestamp, 0).In(loc).Format("2006-01-02")
			aggregate := bson.M{
				"date":     date,
				"location": "Asia/Shanghai",
			}
			_, err = collection.UpdateOne(ctx, bson.M{"_id": position["_id"]}, bson.M{"$set": bson.M{"aggregate": aggregate}})
			if err != nil {
				return
			}

			log.Println("Update transfer_logs success, positionIndex:", position["positionIndex"], ", date:", date)
			fromPositionIndex = position["positionIndex"].(int64)
		}
		log.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
