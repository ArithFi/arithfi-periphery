package model

type Kline struct {
	OpenTime         int64
	Open             float64
	High             float64
	Low              float64
	Close            float64
	Volume           float64
	CloseTime        int64
	QuoteVolume      float64
	NumberOfTrades   int64
	TakerBaseVolume  float64
	TakerQuoteVolume float64
	Ignore           float64
}
