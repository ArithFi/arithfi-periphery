package forex

import (
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/arithfi/arithfi-periphery/model"
	"log"
)

func GetKlines(symbol string, interval string, startTime int64, endTime int64) *[]model.Kline {
	fmt.Println(symbol, interval, startTime, endTime)
	cache := getFromCache(symbol, interval, startTime/1000, endTime/1000)
	if cache != nil {
		fmt.Println("cache hit")
		return cache
	}
	return nil
}

func getFromCache(symbol string, interval string, startTime int64, endTime int64) *[]model.Kline {
	result, _ := mysql.ArithFiDB.Query("select timestamp, open, high, low, close, volume from kline_cache where symbol = ? and resolution = ? and timestamp >= ? and timestamp < ? order by timestamp", symbol, interval, startTime, endTime)

	exchangeInfo := make([]model.Kline, 0)
	for result.Next() {
		var data model.Kline
		err := result.Scan(&data.OpenTime, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			log.Println("Get from cache error")
			return nil
		}
		exchangeInfo = append(exchangeInfo, data)
	}

	return &exchangeInfo
}
