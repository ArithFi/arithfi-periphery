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
		log.Println("futures_positions_retrieve: Scan Futures Tradings from id", lastId)
		actions, err := offchain.GetFuturesTradings(lastId)
		if err != nil {
			log.Println("futures_positions_retrieve: Scan Futures Tradings err", err)
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
						log.Println("futures_positions_retrieve: MARKET_ORDER_FEE", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "MARKET_CLOSE_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"closeFees": action.Fees, "positionStatus": "closed", "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "LIMIT_CANCEL":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "cancelled"}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "LIMIT_EDIT":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"entryPrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
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
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
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
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "MARKET_LIQUIDATION":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionSize": 0, "positionStatus": "closed", "closeFees": 0, "sellValue": 0, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "MARKET_ORDER_ADD":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{
							"margin":      action.Margin,
							"addedMargin": bson.D{{"$subtract", bson.A{action.Margin, "$initialMargin"}}},
						}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "SL_ORDER_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "TP_ORDER_FEE":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"positionStatus": "closed", "closeFees": action.Fees, "sellValue": action.SellValue, "closePrice": action.OrderPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				case "TPSL_EDIT":
					_, err := futuresPositionsCollection.UpdateOne(
						ctx,
						bson.M{"positionIndex": action.PositionIndex},
						bson.M{"$set": bson.M{"stopLossPrice": action.StopLossPrice, "takeProfitPrice": action.TakeProfitPrice}},
					)
					if err != nil {
						log.Println("futures_positions_retrieve:", err)
						break
					}
					log.Println("futures_position_retrieve: success", action.OrderType, action.Id)
					break
				default:
					fmt.Println(action.OrderType, action.Id)
					break
				}
			}
			lastId = actions[len(actions)-1].Id
		} else {
			log.Println("futures_position_retrieve: Scan Futures Tradings empty")
		}
		log.Println("futures_position_retrieve: Sleep 10 seconds")
		time.Sleep(time.Second * 10)
	}
}
