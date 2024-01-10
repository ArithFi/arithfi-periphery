package main

import (
	"context"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/internal/offchain"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func main() {
	lastId := int64(-1)

	ctx := context.TODO()
	futuresPositionsCollection := mongo.MONGODB.Database("off-chain").Collection("futures-positions")

	for {
		log.Println("Scan Futures Tradings from id", lastId)
		actions, err := offchain.GetFuturesTradings(lastId)
		if err != nil {
			log.Println("Scan Futures Tradings err", err)
			continue
		}
		if len(actions) > 0 {
			for _, action := range actions {
				switch action.OrderType {
				case "MARKET_ORDER_FEE":
					_, err := futuresPositionsCollection.InsertOne(ctx, bson.D{
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
						log.Println("MARKET_ORDER_FEE:", err)
						break
					}
					log.Println("MARKET_ORDER_FEE:", action.Id)
					break
				case "MARKET_CLOSE_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"closeFees": action.Fees, "positionStatus": "closed", "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("MARKET_CLOSE_FEE:", err)
						break
					}
					log.Println("MARKET_CLOSE_FEE:", action.Id)
					break
				case "LIMIT_CANCEL":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "cancelled"}},
					)
					if err != nil {
						log.Println("LIMIT_CANCEL:", err)
						break
					}
					log.Println("LIMIT_CANCEL:", action.Id)
					break
				case "LIMIT_EDIT":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"entryPrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("LIMIT_EDIT:", err)
						break
					}
					log.Println("LIMIT_EDIT:", action.Id)
					break
				case "LIMIT_ORDER_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{
							"openPrice":      action.OrderPrice,
							"openFees":       action.Fees,
							"positionStatus": "open",
						}},
					)
					if err != nil {
						log.Println("LIMIT_ORDER_FEE:", err)
						break
					}
					log.Println("LIMIT_ORDER_FEE:", action.Id)
					break
				case "LIMIT_REQUEST":
					_, err := futuresPositionsCollection.InsertOne(ctx, bson.D{
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
						log.Println("LIMIT_REQUEST:", err)
						break
					}
					log.Println("LIMIT_REQUEST:", action.PositionIndex)
					break
				case "MARKET_LIQUIDATION":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionSize": 0, "positionStatus": "closed", "closeFees": 0, "sellValue": 0, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("MARKET_LIQUIDATION:", err)
						break
					}
					log.Println("MARKET_LIQUIDATION:", action.Id)
					break
				case "MARKET_ORDER_ADD":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"margin": action.Margin}},
					)
					if err != nil {
						log.Println("MARKET_ORDER_ADD:", err)
						break
					}
					log.Println("MARKET_ORDER_ADD:", action.Id)
					break
				case "SL_ORDER_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("SL_ORDER_FEE:", err)
						break
					}
					log.Println("SL_ORDER_FEE:", action.Id)
					break
				case "TP_ORDER_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("TP_ORDER_FEE:", err)
						break
					}
					log.Println("TP_ORDER_FEE:", action.Id)
					break
				case "TPSL_EDIT":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"stopLossPrice": action.StopLossPrice, "takeProfitPrice": action.TakeProfitPrice}},
					)
					if err != nil {
						log.Println("TPSL_EDIT:", err)
						break
					}
					log.Println("TPSL_EDIT:", action.Id)
					break
				default:
					fmt.Println(action.OrderType, action.Id)
					break
				}
			}
			lastId = actions[len(actions)-1].Id
		} else {
			log.Println("Scan Futures Tradings empty")
		}
		log.Println("Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
