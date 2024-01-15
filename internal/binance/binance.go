package binance

import (
	"encoding/json"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/arithfi/arithfi-periphery/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	BaseURL   = "https://api.binance.com/api/v3"
	klinesURL = "/klines"
)

func GetKlines(symbol string, interval string, startTime int64, endTime int64, countback int64) *[]model.Kline {
	cache := getFromCache(symbol, interval, startTime/1000, endTime/1000, countback)
	if cache != nil {
		return cache
	}
	body := requestAPI(klinesURL + "?symbol=" + symbol + "&interval=" + interval + "&startTime=" + strconv.FormatInt(startTime, 10) + "&endTime=" + strconv.FormatInt(endTime, 10) + "&limit=" + strconv.FormatInt(countback, 10))
	var arr [][]interface{}
	err := json.Unmarshal(body, &arr)
	if err != nil {
		log.Println(err)
		return nil
	}
	exchangeInfo := make([]model.Kline, len(arr))
	for i, data := range arr {
		exchangeInfo[i] = model.Kline{
			OpenTime: int64(data[0].(float64) / 1000),
			Open:     data[1].(string),
			High:     data[2].(string),
			Low:      data[3].(string),
			Close:    data[4].(string),
			Volume:   data[5].(string),
		}
	}
	go func() {
		if len(exchangeInfo) == 0 {
			return
		}
		err := cacheKlines(&exchangeInfo, interval, symbol)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	return &exchangeInfo
}

func cacheKlines(exchangeInfo *[]model.Kline, resolution string, symbol string) error {
	tx, err := mysql.ArithFiDB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	stmt, err := tx.Prepare(`
		INSERT INTO kline_cache (timestamp, resolution, symbol, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			open = VALUES(open),
			high = VALUES(high),
			low = VALUES(low),
			close = VALUES(close),
			volume = VALUES(volume)
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, data := range *exchangeInfo {
		_, err = stmt.Exec(
			data.OpenTime,
			resolution,
			symbol,
			data.Open,
			data.High,
			data.Low,
			data.Close,
			data.Volume,
		)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func getFromCache(symbol string, interval string, startTime int64, endTime int64, countback int64) *[]model.Kline {
	result, _ := mysql.ArithFiDB.Query("select timestamp, open, high, low, close, volume from kline_cache where symbol = ? and resolution = ? and timestamp >= ? and timestamp < ? order by timestamp limit ?", symbol, interval, startTime, endTime, countback)

	exchangeInfo := make([]model.Kline, 0)
	for result.Next() {
		var data model.Kline
		err := result.Scan(&data.OpenTime, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			log.Println(err)
			return nil
		}
		exchangeInfo = append(exchangeInfo, data)
	}

	if len(exchangeInfo) == 0 {
		return nil
	}

	return &exchangeInfo
}

func requestAPI(endpoint string) []byte {
	resp, err := http.Get(BaseURL + endpoint)
	if err != nil {
		log.Println(err)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	return body
}
