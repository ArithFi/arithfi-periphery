package main

import (
	"context"
	"github.com/arithfi/arithfi-periphery/configs/mongo"
	"github.com/arithfi/arithfi-periphery/internal/offchain"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {
	fromId := int64(0)

	for {
		log.Println("Scan Futures Tradings from id", fromId)
		actions, err := offchain.GetFuturesTradings(fromId)
		if err != nil {
			log.Println("Scan Futures Tradings err", err)
			continue
		}
		if len(actions) == 0 {
			log.Println("Scan Futures Tradings empty")
			break
		}
		fromId = actions[len(actions)-1].PositionIndex + 1
		for _, action := range actions {
			switch action.OrderType {
			case "MARKET_ORDER_FEE":
				handlerMarketOrderFee(&action)
				break
			case "MARKET_CLOSE_FEE":
				handlerMarketCloseFee(&action)
				break
			case "LIMIT_CANCEL":
				handlerLimitCancel(&action)
				break
			case "LIMIT_EDIT":
				handlerLimitEdit(&action)
				break
			case "LIMIT_ORDER_FEE":
				handlerLimitOrderFee(&action)
				break
			case "LIMIT_REQUEST":
				handlerLimitRequest(&action)
				break
			case "MARKET_LIQUIDATION":
				handlerMarketLiquidation(&action)
				break
			case "MARKET_ORDER_ADD":
				handlerMarketOrderAdd(&action)
				break
			case "SL_ORDER_FEE":
				handlerSLOrderFee(&action)
				break
			case "TP_ORDER_FEE":
				handlerTPOrderFee(&action)
				break
			case "TPSL_EDIT":
				handlerTPSLEdit(&action)
				break
			}
		}
	}
}

func handlerMarketOrderFee(action *offchain.Action) {
	ctx := context.TODO()
	collection := mongo.MONGODB.Database("off-chain").Collection("futures-positions")
	_, err := collection.InsertOne(ctx, bson.D{
		{"positionIndex", action.PositionIndex},
		{"product", action.Product},
		{"timeStamp", action.TimeStamp},
		{"leverage", action.Leverage},
		{"positionSize", float64(action.Leverage) * action.Margin},
		{"mode", action.Mode},
		{"direction", action.Direction},
		{"margin", action.Margin},
		{"initialMargin", action.Margin},
		{"walletAddress", action.WalletAddress},
		{"kolAddress", action.KolAddress},
		{"openFees", action.Fees},
	})
	if err != nil {
		return
	}
}

func handlerMarketCloseFee(action *offchain.Action) {
	//log.Println("MarketCloseFee", action)
}

func handlerLimitCancel(action *offchain.Action) {
	//log.Println("LimitCancel", action)
}

func handlerLimitEdit(action *offchain.Action) {
	//log.Println("LimitEdit", action)
}

func handlerLimitOrderFee(action *offchain.Action) {
	//log.Println("LimitOrderFee", action)
}

func handlerLimitRequest(action *offchain.Action) {
	//log.Println("LimitRequest", action)
}

func handlerMarketLiquidation(action *offchain.Action) {
	//log.Println("MarketLiquidation", action)
}

func handlerMarketOrderAdd(action *offchain.Action) {
	//log.Println("MarketOrderAdd", action)
}

func handlerSLOrderFee(action *offchain.Action) {
	//log.Println("SLOrderFee", action)
}

func handlerTPOrderFee(action *offchain.Action) {
	//log.Println("TPOrderFee", action)
}

func handlerTPSLEdit(action *offchain.Action) {
	//log.Println("TPSLEdit", action)
}
