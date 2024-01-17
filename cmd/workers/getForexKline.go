package main

import (
	"encoding/json"
	"fmt"
	"github.com/arithfi/arithfi-periphery/configs/mysql"
	"github.com/arithfi/arithfi-periphery/model"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	BaseURL   = "https://cms.nestfi.net/api/oracle/price"
	klinesURL = "/klines"
)

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

func main() {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick at", time.Now())
			go KlineIntervalWorker("AUDUSD", "500")
			go KlineIntervalWorker("EURUSD", "500")
			go KlineIntervalWorker("USDJPY", "500")
			go KlineIntervalWorker("USDCAD", "500")
			go KlineIntervalWorker("GBPUSD", "500")
		}
	}
}

func KlineIntervalWorker(symbol string, limit string) {
	go GetByInterval(symbol, "1m", limit)
	go GetByInterval(symbol, "5m", limit)
	go GetByInterval(symbol, "15m", limit)
	go GetByInterval(symbol, "30m", limit)
	go GetByInterval(symbol, "1h", limit)
	go GetByInterval(symbol, "1d", limit)
}

func GetByInterval(symbol string, interval string, limit string) *[]model.Kline {
	body := requestAPI(klinesURL + "?symbol=" + symbol + "&interval=" + interval + "&limit=" + limit)
	var arr [][]interface{}
	err := json.Unmarshal(body, &arr)
	if err != nil {
		log.Println("Unmarshal error")
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
			log.Println("Cache error")
			return
		}
	}()
	fmt.Println("Get forex kline: ", symbol, interval, exchangeInfo[0].OpenTime)
	return &exchangeInfo
}
