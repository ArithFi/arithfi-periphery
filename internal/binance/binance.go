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
	if cache != nil && len(*cache) > 0 {
		log.Println("cache hit", symbol, interval, startTime, endTime, len(*cache))
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
				log.Println("api hit", symbol, interval, startTime, endTime)
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

	totalKlines := make([]model.Kline, 0)
	for result.Next() {
		var data model.Kline
		err := result.Scan(&data.OpenTime, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			log.Println("Get from cache error")
			return nil
		}
		totalKlines = append(totalKlines, data)
	}

	if len(totalKlines) == 0 {
		return nil
	}

	var count = (endTime - startTime) / IntervalMap[interval]
	if float64(len(totalKlines)) < float64(int(count))*0.9 {
		return nil
	}
	return &totalKlines
}

var IntervalMap = map[string]int64{
	"1m":  60,
	"3m":  180,
	"5m":  300,
	"15m": 900,
	"30m": 1800,
	"1h":  3600,
	"2h":  7200,
	"4h":  14400,
	"6h":  21600,
	"8h":  28800,
	"12h": 43200,
	"1d":  86400,
	"3d":  259200,
	"1w":  604800,
	"1M":  2592000,
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
