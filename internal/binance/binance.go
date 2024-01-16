package binance

import (
	"encoding/json"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/arithfi/arithfi-periphery/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL   = "https://fapi.binance.com/fapi/v1"
	klinesURL = "/klines"
)

func GetKlines(symbol string, interval string, startTime int64, endTime int64) *[]model.Kline {
	cache := getFromCache(symbol, interval, startTime/1000, endTime/1000)
	if cache != nil {
		return cache
	}

	var from = startTime
	var to = endTime
	var totalKlines []model.Kline

	for {
		body := requestAPI(klinesURL + "?symbol=" + symbol + "&interval=" + interval + "&startTime=" + strconv.FormatInt(from, 10) + "&endTime=" + strconv.FormatInt(to, 10) + "&limit=500")
		var arr [][]interface{}
		err := json.Unmarshal(body, &arr)
		if err != nil {
			log.Println("Unmarshal error")
			return nil
		}
		currKlines := make([]model.Kline, len(arr))
		for i, data := range arr {
			currKlines[i] = model.Kline{
				OpenTime: int64(data[0].(float64) / 1000),
				Open:     data[1].(string),
				High:     data[2].(string),
				Low:      data[3].(string),
				Close:    data[4].(string),
				Volume:   data[5].(string),
			}
		}
		totalKlines = append(totalKlines, currKlines...)
		if len(currKlines) == 500 {
			from = currKlines[len(currKlines)-1].OpenTime + 1
		} else {
			if len(totalKlines) == 0 {
				return nil
			} else {
				go func() {
					err := cacheKlines(&totalKlines, interval, symbol)
					if err != nil {
						log.Println(err)
					}
				}()
				return &totalKlines
			}
		}
		time.Sleep(1 * time.Second)
	}
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
