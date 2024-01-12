package binance

import (
	"encoding/json"
	"github.com/arithfi/arithfi-periphery/model"
	"io"
	"net/http"
	"strconv"
)

const (
	BaseURL   = "https://api.binance.com/api/v3"
	klinesURL = "/klines"
)

func GetKlines(symbol string, interval string, startTime int64, endTime int64) *[]model.Kline {
	body := requestAPI(klinesURL + "?symbol=" + symbol + "&interval=" + interval + "&startTime=" + strconv.FormatInt(startTime, 10) + "&endTime=" + strconv.FormatInt(endTime, 10))
	var arr [][]interface{}
	err := json.Unmarshal(body, &arr)
	if err != nil {
		return nil
	}

	exchangeInfo := make([]model.Kline, len(arr))
	for i, data := range arr {
		exchangeInfo[i] = model.Kline{
			OpenTime:         int64(data[0].(float64) / 1000),
			Open:             data[1].(string),
			High:             data[2].(string),
			Low:              data[3].(string),
			Close:            data[4].(string),
			Volume:           data[5].(string),
			CloseTime:        int64(data[6].(float64) / 1000),
			QuoteVolume:      data[7].(string),
			NumberOfTrades:   int64(data[8].(float64)),
			TakerBaseVolume:  data[9].(string),
			TakerQuoteVolume: data[10].(string),
			Ignore:           data[11].(string),
		}
	}
	return &exchangeInfo
}

func requestAPI(endpoint string) []byte {
	resp, _ := http.Get(BaseURL + endpoint)
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body
}
