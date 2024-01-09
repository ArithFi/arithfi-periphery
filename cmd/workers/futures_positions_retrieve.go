package main

import (
	"github.com/arithfi/arithfi-periphery/internal/offchain"
	"log"
)

type Order struct {
	PositionIndex   int64   `json:"positionIndex"`
	Product         string  `json:"product"`
	TimeStamp       string  `json:"timeStamp"`
	Status          string  `json:"status"`
	Leverage        int64   `json:"leverage"`
	PositionSize    float64 `json:"positionSize"`
	PositionValue   float64 `json:"positionValue"`
	OrderType       string  `json:"orderType"`
	Mode            string  `json:"mode"`
	Direction       string  `json:"direction"`
	Margin          float64 `json:"margin"`
	InitialMargin   float64 `json:"initialMargin"`
	AddedMargin     float64 `json:"addedMargin"`
	UnrealizedPNL   float64 `json:"unrealizedPNL"`
	RealizedPNL     float64 `json:"realizedPNL"`
	OpenPrice       float64 `json:"openPrice"`
	LastPrice       float64 `json:"lastPrice"`
	ClosePrice      float64 `json:"closePrice"`
	MarkPrice       float64 `json:"markPrice"`
	EntryPrice      float64 `json:"entryPrice"`
	UnrealizedROI   float64 `json:"unrealizedROI"`
	RealizedROI     float64 `json:"realizedROI"`
	AssetEquity     float64 `json:"assetEquity"`
	FoundingFees    float64 `json:"foundingFees"`
	OpenFees        float64 `json:"openFees"`
	CloseFees       float64 `json:"closeFees"`
	LiqPrice        float64 `json:"liqPrice"`
	TakeProfitPrice float64 `json:"takeProfitPrice"`
	StopLossPrice   float64 `json:"stopLossPrice"`
	LimitPrice      float64 `json:"limitPrice"`
}

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
			log.Println(action)
		}
	}
}
