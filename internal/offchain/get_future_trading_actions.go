package offchain

import (
	"github.com/arithfi/arithfi-periphery/configs/mysql"
)

type Action struct {
	Id            int64   `json:"id" bson:"id"`
	PositionIndex int64   `json:"positionIndex" bson:"positionIndex"`
	Product       string  `json:"product" bson:"product"`
	TimeStamp     string  `json:"timeStamp" bson:"timeStamp"`
	Leverage      int64   `json:"leverage" bson:"leverage"`
	OrderType     string  `json:"orderType" bson:"orderType"`
	Mode          string  `json:"mode" bson:"mode"`
	Direction     string  `json:"direction" bson:"direction"`
	Margin        float64 `json:"margin" bson:"margin"`
	ClosePrice    float64 `json:"closePrice" bson:"closePrice"`
	WalletAddress string  `json:"walletAddress" bson:"walletAddress"`
	KolAddress    string  `json:"kolAddress" bson:"kolAddress"`
	Fees          float64 `json:"fees" bson:"fees"`
}

// GetFuturesTradings 扫描这个表
func GetFuturesTradings(fromId int64) ([]Action, error) {
	query, err := mysql.MYSQL.Query(`SELECT _id, timeStamp, product, positionIndex, leverage, orderType, mode, direction, margin, volume, sellValue, walletAddress, kolAddress, fees
FROM f_future_trading 
WHERE _id >= ? 
ORDER By _id
LIMIT 10
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
		var timeStamp string
		var leverage int64
		var orderType string
		var mode string
		var direction string
		var margin float64
		var volume float64
		var sellValue float64
		var walletAddress string
		var kolAddress string
		var fees float64

		err := query.Scan(&id, &timeStamp, &product, &positionIndex, &leverage, &orderType, &mode, &direction, &margin, &volume, &sellValue, &walletAddress, &kolAddress, &fees)
		if err != nil {
			return []Action{}, err
		}
		documents = append(documents, Action{
			Id:            id,
			PositionIndex: positionIndex,
			Product:       product,
			TimeStamp:     timeStamp,
			Leverage:      leverage,
			OrderType:     orderType,
			Mode:          mode,
			Direction:     direction,
			Margin:        margin,
			ClosePrice:    sellValue,
			WalletAddress: walletAddress,
			KolAddress:    kolAddress,
			Fees:          fees,
		})
	}
	if err != nil {
		return []Action{}, err
	}
	return documents, nil
}
