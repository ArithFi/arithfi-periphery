package offchain

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"time"
)

type Action struct {
	Id              int64   `json:"id" bson:"id"`
	PositionIndex   int64   `json:"positionIndex" bson:"positionIndex"`
	Product         string  `json:"product" bson:"product"`
	Timestamp       int64   `json:"timestamp" bson:"timestamp"`
	Leverage        int64   `json:"leverage" bson:"leverage"`
	OrderType       string  `json:"orderType" bson:"orderType"`
	OrderPrice      float64 `json:"orderPrice" bson:"orderPrice"`
	Mode            string  `json:"mode" bson:"mode"`
	Direction       string  `json:"direction" bson:"direction"`
	Margin          float64 `json:"margin" bson:"margin"`
	ClosePrice      float64 `json:"closePrice" bson:"closePrice"`
	WalletAddress   string  `json:"walletAddress" bson:"walletAddress"`
	KolAddress      string  `json:"kolAddress" bson:"kolAddress"`
	Fees            float64 `json:"fees" bson:"fees"`
	StopLossPrice   float64 `json:"stopLossPrice" bson:"stopLossPrice"`
	TakeProfitPrice float64 `json:"takeProfitPrice" bson:"takeProfitPrice"`
	SellValue       float64 `json:"sellValue" bson:"sellValue"`
}

// GetFuturesTradings 扫描这个表
func GetFuturesTradings(fromId int64) ([]Action, error) {
	query, err := mysql.MYSQL.Query(`SELECT _id, timestamp, product, positionIndex, leverage, orderType, orderPrice, mode, direction, margin, volume, sellValue, walletAddress, kolAddress, fees, stopLossPrice, takeProfitPrice
FROM f_future_trading 
WHERE _id > ? 
ORDER By _id
LIMIT 1000
`, fromId)
	if err != nil {
		return []Action{}, err
	}
	defer query.Close()
	var documents []Action
	for query.Next() {
		var id int64
		var product string
		var positionIndex int64
		var timeStampStr string
		var leverage int64
		var orderType string
		var orderPrice float64
		var mode string
		var direction string
		var margin float64
		var volume float64
		var sellValue float64
		var walletAddress string
		var kolAddress string
		var fees float64
		var stopLossPrice float64
		var takeProfitPrice float64

		err := query.Scan(&id, &timeStampStr, &product, &positionIndex, &leverage, &orderType, &orderPrice, &mode, &direction, &margin, &volume, &sellValue, &walletAddress, &kolAddress, &fees, &stopLossPrice, &takeProfitPrice)
		if err != nil {
			return []Action{}, err
		}

		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", timeStampStr, location)
		if direction == "1" {
			direction = "LONG"
		} else {
			direction = "SHORT"
		}
		documents = append(documents, Action{
			Id:              id,
			PositionIndex:   positionIndex,
			Product:         product,
			Timestamp:       timestamp.Unix(),
			Leverage:        leverage,
			OrderType:       orderType,
			OrderPrice:      orderPrice,
			Mode:            mode,
			Direction:       direction,
			Margin:          margin,
			ClosePrice:      sellValue,
			WalletAddress:   walletAddress,
			KolAddress:      kolAddress,
			Fees:            fees,
			StopLossPrice:   stopLossPrice,
			TakeProfitPrice: takeProfitPrice,
			SellValue:       sellValue,
		})
	}
	if err != nil {
		return []Action{}, err
	}
	return documents, nil
}
