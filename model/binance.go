package model

type Kline struct {
	OpenTime         int64
	Open             string
	High             string
	Low              string
	Close            string
	Volume           string
	CloseTime        int64
	QuoteVolume      string
	NumberOfTrades   int64
	TakerBaseVolume  string
	TakerQuoteVolume string
	Ignore           string
}
