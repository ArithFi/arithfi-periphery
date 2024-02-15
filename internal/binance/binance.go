package binance

import (
	"context"
	"encoding/json"
	"github.com/arithfi/arithfi-periphery/configs/cache"
	"github.com/arithfi/arithfi-periphery/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL   = "https://88i954i3v1.execute-api.ap-northeast-1.amazonaws.com/prod"
	klinesURL = "/klines"
)

func GetKlines(symbol string, interval string, startTime int64, endTime int64) *[]model.Kline {
	var from = startTime
	var to = endTime
	var totalKlines []model.Kline
	ctx := context.Background()

	for {
		var uri = klinesURL + "?symbol=" + symbol + "&interval=" + interval + "&startTime=" + strconv.FormatInt(from, 10) + "&endTime=" + strconv.FormatInt(to, 10) + "&limit=500"
		cacheArrayCmd := cache.CACHE.Get(ctx, uri)

		if cacheArrayCmd.Err() == nil {
			cacheArrayStr := cacheArrayCmd.Val()

			if cacheArrayStr != "" {
				var cacheKlines []model.Kline
				cacheArray := []byte(cacheArrayStr)
				err := json.Unmarshal(cacheArray, &cacheKlines)
				if err != nil {
					log.Println("Unmarshal error")
				} else {
					log.Println("Get from cache: ", uri)
					return &cacheKlines
				}
			}
		}

		body := requestAPI(uri)
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
				totalKlinesJSON, err := json.Marshal(totalKlines)
				if err != nil {
					log.Fatalf("Failed to marshal klines: %v", err)
				}
				totalKlinesStr := string(totalKlinesJSON)
				_, err = cache.CACHE.Set(ctx, uri, totalKlinesStr, 2*time.Second).Result()
				if err != nil {
					log.Fatalf("Failed to set klines in Redis: %v", err)
				}
				log.Println("Set to cache: ", uri)
				return &totalKlines
			}
		}
		time.Sleep(1 * time.Second)
	}
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
