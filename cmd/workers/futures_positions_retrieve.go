package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/internal/offchain"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {
	lastId := int64(-1)

	ctx := context.TODO()
	collection := mongo.MONGODB.Database("off-chain").Collection("futures-positions")

	for {
		log.Println("Scan Futures Tradings from id", lastId)
		actions, err := offchain.GetFuturesTradings(lastId)
		if err != nil {
			log.Println("Scan Futures Tradings err", err)
			continue
		}
		if len(actions) == 0 {
			log.Println("Scan Futures Tradings empty")
			break
		}
		lastId = actions[len(actions)-1].PositionIndex
		for _, action := range actions {
			switch action.OrderType {
			case "MARKET_ORDER_FEE":
				_, err := collection.InsertOne(ctx, bson.D{
					{"positionIndex", action.PositionIndex},
					{"product", action.Product},
					{"positionStatus", "open"},
					{"timestamp", action.Timestamp},
					{"leverage", action.Leverage},
					{"positionSize", float64(action.Leverage) * action.Margin},
					{"mode", action.Mode},
					{"direction", action.Direction},
					{"margin", action.Margin},
					{"initialMargin", action.Margin},
					{"walletAddress", action.WalletAddress},
					{"kolAddress", action.KolAddress},
					{"openFees", action.Fees},
					{"openPrice", action.OrderPrice},
				})
				if err != nil {
					break
				}
				log.Println("MARKET_ORDER_FEE:", action.PositionIndex)
				break
			case "MARKET_CLOSE_FEE":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"closeFees": action.Fees, "positionStatus": "closed", "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
				)
				if err != nil {
					break
				}
				log.Println("MARKET_CLOSE_FEE:", action.PositionIndex)
				break
			case "LIMIT_CANCEL":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"positionStatus": "cancelled"}},
				)
				if err != nil {
					break
				}
				log.Println("LIMIT_CANCEL:", action.PositionIndex)
				break
			case "LIMIT_EDIT":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"entryPrice": action.OrderPrice}},
				)
				if err != nil {
					break
				}
				log.Println("LIMIT_EDIT:", action.PositionIndex)
				break
			case "LIMIT_ORDER_FEE":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{
						"openPrice":      action.OrderPrice,
						"openFees":       action.Fees,
						"positionStatus": "open",
					}},
				)
				if err != nil {
					break
				}
				log.Println("LIMIT_ORDER_FEE:", action.PositionIndex)
				break
			case "LIMIT_REQUEST":
				_, err := collection.InsertOne(ctx, bson.D{
					{"positionIndex", action.PositionIndex},
					{"product", action.Product},
					{"positionStatus", "pending"},
					{"timestamp", action.Timestamp},
					{"leverage", action.Leverage},
					{"positionSize", float64(action.Leverage) * action.Margin},
					{"mode", action.Mode},
					{"direction", action.Direction},
					{"margin", action.Margin},
					{"initialMargin", action.Margin},
					{"walletAddress", action.WalletAddress},
					{"kolAddress", action.KolAddress},
					{"entryPrice", action.OrderPrice},
				})
				if err != nil {
					break
				}
				log.Println("LIMIT_REQUEST:", action.PositionIndex)
				break
			case "MARKET_LIQUIDATION":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"positionSize": 0, "positionStatus": "closed", "closeFees": 0, "sellValue": 0, "closePrice": action.OrderPrice}},
				)
				if err != nil {
					break
				}
				log.Println("MARKET_LIQUIDATION:", action.PositionIndex)
				break
			case "MARKET_ORDER_ADD":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"margin": action.Margin}},
				)
				if err != nil {
					break
				}
				log.Println("MARKET_ORDER_ADD:", action.PositionIndex)
				break
			case "SL_ORDER_FEE":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
				)
				if err != nil {
					break
				}
				log.Println("SL_ORDER_FEE:", action.PositionIndex)
				break
			case "TP_ORDER_FEE":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
				)
				if err != nil {
					break
				}
				log.Println("TP_ORDER_FEE:", action.PositionIndex)
				break
			case "TPSL_EDIT":
				_, err := collection.UpdateOne(
					ctx,
					bson.M{"positionIndex": action.PositionIndex},
					bson.M{"$set": bson.M{"stopLossPrice": action.StopLossPrice, "takeProfitPrice": action.TakeProfitPrice}},
				)
				if err != nil {
					break
				}
				log.Println("TPSL_EDIT:", action.PositionIndex)
				break
			}
		}
	}
}
